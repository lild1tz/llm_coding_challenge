from fastapi import FastAPI
from app.utils import convert_ogg_to_mp3, get_embedding, base64_to_dataurl
from app.prompts import fewshot_text_prompt, message_photo_prompt, message_audio_prompt, change_table_prompt
from app.models import Table, ClassifyMessageOutput, InputMessage, InputPhoto, InputAudio, OutputAudio, InputChangeTable
from app.config import Config
from langchain_openai import ChatOpenAI
from transformers import AutoTokenizer, AutoConfig
import onnxruntime as ort
from joblib import load
from openai import AsyncOpenAI
from pathlib import Path
import base64
import io
import warnings
warnings.filterwarnings("ignore")

config = Config()

app = FastAPI(
    title="Apollo API",
    description="API для обработки сообщений, классификации, обработки фото и аудио, а также изменения таблиц.",
    version="1.0.0"
)

llm = ChatOpenAI(
   model_name=config.OPENAI_MODEL,
   openai_api_key=config.OPENAI_API_KEY,
   base_url=config.OPENAI_BASE_URL,
   temperature=0,
   max_tokens=10240
).with_structured_output(Table)

transcriber = AsyncOpenAI(
    api_key=config.OPENAI_API_KEY,
    base_url=config.OPENAI_BASE_URL 
)

base_dir = Path(__file__).resolve().parent
onnx_model_path = base_dir / "source" / "onnx_model" / "model.onnx"
tokenizer_dir = base_dir / "source" / "onnx_model"
logreg_model_path = base_dir / "source" / "logreg_model.joblib"
session = ort.InferenceSession(onnx_model_path.as_posix())
tokenizer_config = AutoConfig.from_pretrained(tokenizer_dir.as_posix(), local_files_only=True)
tokenizer = AutoTokenizer.from_pretrained(tokenizer_dir.as_posix(), config=tokenizer_config, local_files_only=True)
clf = load(logreg_model_path)

@app.post("/process_message", response_model=Table, summary="Обработка сообщения", description="Обрабатывает текстовое сообщение и возвращает таблицу.")
async def process_message(input: InputMessage) -> Table:
    massage = input.message
    prompt = fewshot_text_prompt.format_prompt(input=massage)
    table = await llm.ainvoke(prompt)
    return table

@app.post("/classify_message", response_model=ClassifyMessageOutput, summary="Классификация сообщения", description="Классифицирует сообщение и возвращает вероятность и предсказанный класс.")
async def classify_message(input: InputMessage) -> ClassifyMessageOutput:
    massage = input.message
    embedding = get_embedding(massage, tokenizer, session)
    prediction = clf.predict_proba([embedding])[0]
    predicted_class = int(prediction.argmax())
    probability = round(prediction[1], 2)
    return ClassifyMessageOutput(probability=probability, prediction=predicted_class)

@app.post("/process_photo", response_model=Table, summary="Обработка фото", description="Обрабатывает фото и возвращает таблицу.")
async def process_photo(input: InputPhoto) -> Table:
    base64_str = input.photo
    file_type = input.type
    dataurl = base64_to_dataurl(base64_str, file_type)
    prompt = message_photo_prompt(dataurl)
    table = await llm.ainvoke(prompt)
    return table

@app.post("/transcribe_audio", response_model=OutputAudio, summary="Транскрипция аудио", description="Транскрибирует аудиофайл и возвращает текст.")
async def transcribe_audio(input: InputAudio) -> OutputAudio:
    base64_str = input.audio
    file_type = input.type.lower()
    audio_data = base64.b64decode(base64_str)

    audio_file = io.BytesIO(audio_data)
    audio_file.name = f"audio.{file_type}"

    if file_type == "ogg":
        audio_file = convert_ogg_to_mp3(audio_data)
        audio_file.name = f"audio.mp3"

    transcription = await transcriber.audio.transcriptions.create(
        model=config.TRANSCRIBER_MODEL,
        file=audio_file
    )
    
    return OutputAudio(text=transcription.text)

@app.post("/change_table", response_model=Table, summary="Изменение таблицы", description="Изменяет таблицу на основе входных данных и возвращает обновленную таблицу.")
async def change_table(input: InputChangeTable) -> Table:
    table = input.table
    message = input.message
    prompt = change_table_prompt.format_prompt(input=table, message=message)
    table = await llm.ainvoke(prompt)
    return table

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
    #uvicorn app.main:app --reload
