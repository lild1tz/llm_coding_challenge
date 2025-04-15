from fastapi import FastAPI
from app.utils import get_embedding, base64_to_dataurl
from app.prompts import fewshot_text_prompt, message_photo_prompt
from app.models import Table, ClassifyMessageOutput, InputMessage, InputPhoto, InputAudio, OutputAudio
from app.config import Config
from langchain_openai import ChatOpenAI
from transformers import AutoTokenizer, AutoConfig
import onnxruntime as ort
from joblib import load
from dotenv import load_dotenv
from pathlib import Path
import warnings

warnings.filterwarnings("ignore")
config = Config()

app = FastAPI()
print(config.OPENAI_MODEL, config.OPENAI_API_KEY, config.OPENAI_BASE_URL)

llm = ChatOpenAI(
   model=config.OPENAI_MODEL,
   openai_api_key=config.OPENAI_API_KEY,
   base_url=config.OPENAI_BASE_URL,
#    temperature=0,
   max_tokens=10240
).with_structured_output(Table)

transcriber = ChatOpenAI(
    model=config.TRANSCRIBER_MODEL,
    openai_api_key= config.OPENAI_API_KEY,
    base_url=config.OPENAI_BASE_URL
)


base_dir = Path(__file__).resolve().parent
onnx_model_path = base_dir / "source" / "onnx_model" / "model.onnx"
tokenizer_dir = base_dir / "source" / "onnx_model"
logreg_model_path = base_dir / "source" / "logreg_model.joblib"
session = ort.InferenceSession(onnx_model_path.as_posix())
config = AutoConfig.from_pretrained(tokenizer_dir.as_posix(), local_files_only=True)
tokenizer = AutoTokenizer.from_pretrained(tokenizer_dir.as_posix(), config=config, local_files_only=True)
clf = load(logreg_model_path)

@app.post("/process_message")
async def process_message(input: InputMessage) -> Table:
    massage = input.message
    prompt = fewshot_text_prompt.format_prompt(input=massage)
    table = await llm.ainvoke(prompt) #
    return table

@app.post("/classify_message")
async def classify_message(input: InputMessage) -> ClassifyMessageOutput:
    massage = input.message
    embedding = get_embedding(massage, tokenizer, session)
    prediction = clf.predict_proba([embedding])[0]
    predicted_class = int(prediction.argmax())
    probability = round(prediction[1], 2)
    return ClassifyMessageOutput(probability=probability, prediction=predicted_class)

@app.post("/process_photo")
async def process_photo(input: InputPhoto) -> Table:
    base64_str = input.photo
    file_type = input.type
    dataurl = base64_to_dataurl(base64_str, file_type)
    prompt = message_photo_prompt(dataurl)
    table = await llm.ainvoke(prompt) 
    return table
   
@app.post("/transcribe_audio")
async def transcribe_audio(input: InputAudio) -> OutputAudio:
    base64_str = input.audio
    file_type = input.type
    dataurl = base64_to_dataurl(base64_str, file_type)
    prompt = transcriber.invoke(dataurl)
    return OutputAudio(text=prompt)
    
if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
    #uvicorn app.main:app --reload




