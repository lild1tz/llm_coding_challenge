from fastapi import FastAPI
from app.utils import get_embedding
from app.prompts import fewshot_prompt
from app.models import ProcessMessageOutput, ClassifyMessageOutput, InputMessage
from app.config import Config
from langchain_openai import ChatOpenAI
from transformers import AutoTokenizer, AutoConfig
import onnxruntime as ort
from joblib import load
from pathlib import Path

import warnings
warnings.filterwarnings("ignore")

app = FastAPI()

llm = ChatOpenAI(
    model_name=Config.OPENAI_MODEL,
    openai_api_key=Config.OPENAI_API_KEY,
    base_url=Config.OPENAI_BASE_URL
).with_structured_output(ProcessMessageOutput)

base_dir = Path(__file__).resolve().parent
onnx_model_path = base_dir / "source" / "onnx_model" / "model.onnx"
tokenizer_dir = base_dir / "source" / "onnx_model"
logreg_model_path = base_dir / "source" / "logreg_model.joblib"
session = ort.InferenceSession(onnx_model_path.as_posix())
config = AutoConfig.from_pretrained(tokenizer_dir.as_posix(), local_files_only=True)
tokenizer = AutoTokenizer.from_pretrained(tokenizer_dir.as_posix(), config=config, local_files_only=True)
clf = load(logreg_model_path)

@app.post("/process_message")
async def process_message(input: InputMessage) -> ProcessMessageOutput:
    massage = input.message
    prompt = fewshot_prompt.format_prompt(input=massage)
    table = await llm.ainvoke(prompt)
    return table

@app.post("/classify_message")
async def classify_message(input: InputMessage) -> ClassifyMessageOutput:
    massage = input.message
    embedding = get_embedding(massage, tokenizer, session)
    prediction = clf.predict_proba([embedding])[0]
    predicted_class = int(prediction.argmax())
    probability = round(prediction[1], 2)
    return ClassifyMessageOutput(probability=probability, prediction=predicted_class)

if __name__ == "__main__":
    import uvicorn
    print(onnx_model_path)
    uvicorn.run(app, host="0.0.0.0", port=8000)
    #uvicorn app.main:app --reload


