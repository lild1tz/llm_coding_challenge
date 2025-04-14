from dotenv import load_dotenv
import os

load_dotenv()

class Config:
    OPENAI_API_KEY = os.getenv('OPENAI_API_KEY')
    OPENAI_MODEL = os.getenv('OPENAI_MODEL')
    OPENAI_BASE_URL = os.getenv('OPENAI_BASE_URL')
    LANGSMITH_TRACING = os.getenv('LANGSMITH_TRACING')
    LANGSMITH_ENDPOINT = os.getenv('LANGSMITH_ENDPOINT')
    LANGSMITH_API_KEY = os.getenv('LANGSMITH_API_KEY')
    LANGSMITH_PROJECT = os.getenv('LANGSMITH_PROJECT')
    BERT_MODEL = os.getenv('BERT_MODEL')