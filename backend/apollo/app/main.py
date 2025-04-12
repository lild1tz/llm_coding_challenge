from fastapi import FastAPI
from app.prompts import fewshot_prompt
from app.models import Output, Input
from app.config import Config
from langchain_openai import ChatOpenAI
import logging

app = FastAPI()

llm = ChatOpenAI(
    model_name=Config.OPENAI_MODEL,
    openai_api_key=Config.OPENAI_API_KEY,
    base_url=Config.OPENAI_BASE_URL,
).with_structured_output(Output)

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


# massage = "Пахота зяби под мн тр\nПо Пу 26/488\nОтд 12 26/221\n\nПредп культ под оз пш\nПо Пу 215/1015\nОтд 12 128/317\nОтд 16 123/529\n\n2-е диск сах св под оз пш\nПо Пу 22/627\nОтд 11 22/217\n\n2-е диск сои под оз пш\nПо Пу 45/1907\nОтд 12 45/299"
# formatted_prompt = fewshot_prompt.format(input=massage)
# print(formatted_prompt)

@app.post("/process_message")
async def process_message(input: Input) -> Output:
    massage = input.message
    logger.info(f"Received message: {massage}")
    prompt = fewshot_prompt.format_prompt(input=massage)
    logger.info(f"Prompt: {prompt.text}")
    table = await llm.ainvoke(prompt)
    logger.info(f"Table: {table}")
    return table

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
#     #uvicorn app.main:app --reload


