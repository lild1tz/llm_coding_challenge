import requests
import pandas as pd
from sklearn.metrics import accuracy_score
import json
from tqdm import tqdm
from app.config import Config
# Чтение данных из Excel-файла
excel_file = "benchmark_testing.xlsx"
data_df = pd.read_excel(excel_file)

# URL вашего API для тестирования
config = Config()
print(config.OPENAI_MODEL, config.OPENAI_API_KEY, config.OPENAI_BASE_URL)
API_URL = "http://localhost:8000/process_message"
MODEL_NAME = config.OPENAI_MODEL


# Сбор результатов
results = []
model_responses = []

for _, row in tqdm(data_df.iterrows(), total=len(data_df)):
    message = row['message']
    expected_output = json.loads(row['output'])

    response = requests.post(API_URL, json={"message": message})
    output = response.json()
    model_responses.append({"message": message, "model_output": output})

    for pred, exp in zip(output["table"], expected_output["table"]):
        result = {
            field: int(str(pred[field]).strip() == str(exp[field]).strip()) if pred[field] is not None else None
            for field in pred
        }
        results.append(result)

# Преобразование в DataFrame для анализа
results_df = pd.DataFrame(results)

# Расчёт accuracy по каждому полю
accuracies = {}
for column in results_df.columns:
    valid_entries = results_df[column].dropna()
    if not valid_entries.empty:
        accuracies[column] = accuracy_score([1]*len(valid_entries), valid_entries)

# Вывод результатов
print(f"Accuracy по каждому полю {MODEL_NAME}:")
for field, acc in accuracies.items():
    print(f"{field}: {acc:.2f}")

# Сохранение результатов в CSV
results_df.to_csv(f"benchmark_result/benchmark_results_{MODEL_NAME}.csv", index=False)

# Сохранение accuracy в JSON
with open(f"benchmark_result/accuracy_results_{MODEL_NAME}.json", "w") as f:
    json.dump(accuracies, f, ensure_ascii=False, indent=2)

# Сохранение ответов модели в JSON
with open(f"benchmark_result/model_responses_{MODEL_NAME}.json", "w") as f:
    json.dump(model_responses, f, ensure_ascii=False, indent=2)

print(f"Benchmark завершён, результаты сохранены в benchmark_results_{MODEL_NAME}.csv, accuracy_results_{MODEL_NAME}.json и model_responses_{MODEL_NAME}.json")