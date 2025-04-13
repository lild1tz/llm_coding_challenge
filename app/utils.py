import json

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

def ocr_photo(model, photo_base64: str) -> str:
    pass
    
