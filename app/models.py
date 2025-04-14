from pydantic import BaseModel, Field
from typing import List, Union, Optional, Literal

class TableRow(BaseModel):
    date: Optional[str] = Field(None, description="Дата выполнения операции, если представлена в сообщении, в формате ДД.ММ или ДД.ММ.ГГГГ")
    division: str = Field(..., description="Подразделение, выполнявшее работу (например, АОР, Юг, Мир и т.д.)") #
    operation: str = Field(..., description="Наименование выполненной полевой операции (например, Пахота, Дискование и т.д.)") #
    culture: str = Field(..., description="Сельскохозяйственная культура, к которой относится операция") #
    per_day: Union[int, float, None] = Field(None, description="Объём выполненных работ за день (в гектарах)") #
    per_operation: Union[int, float, None] = Field(None, description="Накопленный объём работ с начала операции (в гектарах)") #
    val_day: Union[int, float, None] = Field(None, description="Валовый сбор за день (в центнерах), если применимо")
    val_beginning: Union[int, float, None] = Field(None, description="Суммарный валовый сбор с начала операции (в центнерах), если применимо")

class Table(BaseModel):
    table: List[TableRow]

class InputMessage(BaseModel):
    message: str = Field(..., description="Сообщение от пользователя")

class ClassifyMessageOutput(BaseModel):
    probability: float = Field(..., description="Вероятность того, что сообщение относится к классу 'Операция'")
    prediction: Literal[0, 1] = Field(..., description="Предсказанный класс сообщения (0 - не операция, 1 - операция)")
    
class InputPhoto(BaseModel):
    photo: str = Field(..., description="Фотография в формате base64")
    type: Literal["png", "jpeg", "jpg"] = Field(..., description="Тип фотографии")
