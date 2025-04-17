import io
import json
from pydub import AudioSegment

def safe_json(text):
    return json.dumps(text, ensure_ascii=False, indent=2).replace("{", "{{").replace("}", "}}")

def get_embedding(text: str, tokenizer, session):
    inputs = tokenizer(text, return_tensors="np", padding=True, truncation=True)
    ort_inputs = {
        "input_ids": inputs["input_ids"],
        "attention_mask": inputs["attention_mask"]
    }
    ort_outputs = session.run(None, ort_inputs)
    embedding = ort_outputs[0].mean(axis=1)[0]
    return embedding

def base64_to_dataurl(base64_str: str, file_type: str) -> str:
    return f"data:image/{file_type};base64,{base64_str}"

def convert_ogg_to_mp3(ogg_bytes: bytes) -> bytes:
    audio = AudioSegment.from_ogg(io.BytesIO(ogg_bytes))
    
    audio = audio.set_channels(1)
    audio = audio.set_frame_rate(16000)
    
    mp3_data = io.BytesIO()
    audio.export(
        mp3_data,
        format="mp3",
        codec="libmp3lame",
        parameters=["-b:a", "64k", "-aq", "2"]
    )
    
    return mp3_data.getvalue()
