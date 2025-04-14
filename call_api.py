import requests
import base64

def send_photo_for_processing(photo_path: str, file_type: str):
    with open(photo_path, "rb") as image_file:
        base64_str = base64.b64encode(image_file.read()).decode('utf-8')
    
    payload = {
        "photo": base64_str,
        "type": file_type
    }
    
    response = requests.post("http://127.0.0.1:8000/process_photo", json=payload)
    
    if response.status_code == 200:
        return response.json()
    else:
        response.raise_for_status()

if __name__ == "__main__":
    photo_path = "IMG-20250113-WA0008.jpg"
    result = send_photo_for_processing(photo_path, "jpg")
    print(result)

