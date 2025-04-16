from dotenv import load_dotenv
import os

load_dotenv()

class Config:
    def __init__(self):
        self.OPENAI_API_KEY = os.getenv('OPENAI_API_KEY')
        self.OPENAI_MODEL = os.getenv('OPENAI_MODEL')
        self.OPENAI_BASE_URL = os.getenv('OPENAI_BASE_URL')
        self.LANGSMITH_TRACING = os.getenv('LANGSMITH_TRACING')
        self.LANGSMITH_ENDPOINT = os.getenv('LANGSMITH_ENDPOINT')
        self.LANGSMITH_API_KEY = os.getenv('LANGSMITH_API_KEY')
        self.LANGSMITH_PROJECT = os.getenv('LANGSMITH_PROJECT')
        self.BERT_MODEL = os.getenv('BERT_MODEL')
        self.TRANSCRIBER_MODEL = os.getenv('TRANSCRIBER_MODEL')