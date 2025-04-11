import json

def safe_json(text):
    return json.dumps(text, ensure_ascii=False, indent=2).replace("{", "{{").replace("}", "}}")

