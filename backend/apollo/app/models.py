from pydantic import BaseModel, Field
from typing import List, Union, Optional

class TableRow(BaseModel):
    date: Optional[str] = Field(None, description="Дата выполнения операции, если представлена в сообщении, в формате ДД.ММ")
    division: str = Field(..., description="Подразделение, выполнявшее работу (например, АОР, Юг, Мир и т.д.)")
    operation: str = Field(..., description="Наименование выполненной полевой операции (например, Пахота, Дискование и т.д.)")
    culture: str = Field(..., description="Сельскохозяйственная культура, к которой относится операция")
    per_day: Union[int, float] = Field(..., description="Объём выполненных работ за день (в гектарах)")
    per_operation: Union[int, float] = Field(..., description="Накопленный объём работ с начала операции (в гектарах)")
    val_day: Union[int, float] = Field(None, description="Валовый сбор за день (в центнерах), если применимо")
    val_beginning: Union[int, float] = Field(None, description="Суммарный валовый сбор с начала операции (в центнерах), если применимо")

class Output(BaseModel):
    table: List[TableRow]

class Input(BaseModel):
    message: str = Field(..., description="Сообщение от пользователя")



