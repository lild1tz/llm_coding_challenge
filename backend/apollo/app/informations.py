from app.utils import safe_json

cultures = [
      "Вика+Тритикале",
      "Горох на зерно",
      "Горох товарный",
      "Гуар",
      "Конопля",
      "Кориандр",
      "Кукуруза кормовая",
      "Кукуруза семенная",
      "Кукуруза товарная",
      "Люцерна",
      "Многолетние злаковые травы",
      "Многолетние травы прошлых лет",
      "Многолетние травы текущего года",
      "Овес",
      "Подсолнечник кондитерский",
      "Подсолнечник семенной",
      "Подсолнечник товарный",
      "Просо",
      "Пшеница озимая на зеленый корм",
      "Пшеница озимая семенная",
      "Пшеница озимая товарная",
      "Рапс озимый",
      "Рапс яровой",
      "Свекла сахарная",
      "Сорго",
      "Сорго кормовой",
      "Сорго-суданковый гибрид",
      "Соя семенная",
      "Соя товарная",
      "Чистый пар",
      "Чумиза",
      "Ячмень озимый",
      "Ячмень озимый семенной"
    ]


units = [
    { "division": "АОР", "pu": "Кавказ", "department": "18" },
    { "division": "АОР", "pu": "Кавказ", "department": "19" },
    { "division": "АОР", "pu": "Север", "department": "3" },
    { "division": "АОР", "pu": "Север", "department": "7" },
    { "division": "АОР", "pu": "Север", "department": "10" },
    { "division": "АОР", "pu": "Север", "department": "20" },
    { "division": "АОР", "pu": "Центр", "department": "1" },
    { "division": "АОР", "pu": "Центр", "department": "4" },
    { "division": "АОР", "pu": "Центр", "department": "5" },
    { "division": "АОР", "pu": "Центр", "department": "6" },
    { "division": "АОР", "pu": "Центр", "department": "9" },
    { "division": "АОР", "pu": "Юг", "department": "11" },
    { "division": "АОР", "pu": "Юг", "department": "12" },
    { "division": "АОР", "pu": "Юг", "department": "16" },
    { "division": "АОР", "pu": "Юг", "department": "17" },
    { "division": "ТСК", "pu": "Нет ПУ", "department": "Нет отделения" },
    { "division": "АО Кропоткинское", "pu": "Нет ПУ", "department": "Нет отделения" },
    { "division": "Восход", "pu": "Нет ПУ", "department": "Нет отделения" },
    { "division": "Колхоз Прогресс", "pu": "Нет ПУ", "department": "Нет отделения" },
    { "division": "Мир", "pu": "Нет ПУ", "department": "Нет отделения" },
    { "division": "СП Коломейцево", "pu": "Нет ПУ", "department": "Нет отделения" }
  ]
  
operations = [
      {
        "name": "1-я междурядная культивация",
        "note": "На всех культурах кроме пшеницы, ячменя"
      },
      {
        "name": "2-я междурядная культивация",
        "note": "На всех культурах кроме пшеницы, ячменя"
      },
      {
        "name": "Боронование довсходовое",
        "note": ""
      },
      {
        "name": "Внесение минеральных удобрений",
        "note": ""
      },
      {
        "name": "Выравнивание зяби",
        "note": ""
      },
      {
        "name": "2-е Выравнивание зяби",
        "note": ""
      },
      {
        "name": "Гербицидная обработка",
        "note": "На свекле их 4 шт, на остальных культурах 1"
      },
      {
        "name": "1 Гербицидная обработка",
        "note": "На свекле их 4 шт, на остальных культурах 1"
      },
      {
        "name": "2 Гербицидная обработка",
        "note": "На свекле их 4 шт, на остальных культурах 1"
      },
      {
        "name": "3 Гербицидная обработка",
        "note": "На свекле их 4 шт, на остальных культурах 1"
      },
      {
        "name": "4 Гербицидная обработка",
        "note": "На свекле их 4 шт, на остальных культурах 1"
      },
      {
        "name": "Дискование",
        "note": ""
      },
      {
        "name": "Дискование 2-е",
        "note": ""
      },
      {
        "name": "Инсектицидная обработка",
        "note": ""
      },
      {
        "name": "Культивация",
        "note": ""
      },
      {
        "name": "Пахота",
        "note": ""
      },
      {
        "name": "Подкормка",
        "note": ""
      },
      {
        "name": "Предпосевная культивация",
        "note": ""
      },
      {
        "name": "Прикатывание посевов",
        "note": ""
      },
      {
        "name": "Сев",
        "note": ""
      },
      {
        "name": "Сплошная культивация",
        "note": ""
      },
      {
        "name": "Уборка",
        "note": ""
      },
      {
        "name": "Функицидная обработка",
        "note": ""
      },
      {
        "name": "Чизлевание",
        "note": ""
      }
    ]

examples = [
  {
    "input": "Пахота зяби под мн тр\nПо Пу 26/488\nОтд 12 26/221\n\nПредп культ под оз пш\nПо Пу 215/1015\nОтд 12 128/317\nОтд 16 123/529\n\n2-е диск сах св под оз пш\nПо Пу 22/627\nОтд 11 22/217\n\n2-е диск сои под оз пш\nПо Пу 45/1907\nОтд 12 45/299",
    "output": safe_json([
      {
        "date": "",
        "division": "АОР",
        "operation": "Пахота",
        "culture": "Многолетние травы",
        "per_day": 26,
        "per_operation": 488,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Предпосевная культивация",
        "culture": "Пшеница озимая товарная",
        "per_day": 215,
        "per_operation": 1015,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Дискование 2-е",
        "culture": "Пшеница озимая товарная",
        "per_day": 22,
        "per_operation": 627,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Дискование 2-е",
        "culture": "Пшеница озимая товарная",
        "per_day": 45,
        "per_operation": 1907,
        "val_day": "",
        "val_beginning": ""
      }
    ]),
  },
  {
    "input": "Пахота зяби под сою\nПо ПУ 7/1402\nОтд 17 7/141\n\nВырав-ие зяби под кук/силос\nПо ПУ 16/16\nОтд 12 16/16\n\nВырав-ие зяби под сах/свёклу\nПо ПУ 67/912\nОтд 12 67/376\n\n2-ое диск-ие сах/свёкла\nПо ПУ 59/1041\nОтд 17 59/349",
    "output": safe_json([
      {
        "date": "",
        "division": "АОР",
        "operation": "Пахота",
        "culture": "Соя",
        "per_day": 7,
        "per_operation": 1402,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Выравнивание зяби",
        "culture": "Кукуруза кормовая",
        "per_day": 16,
        "per_operation": 16,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Выравнивание зяби",
        "culture": "Свекла сахарная",
        "per_day": 67,
        "per_operation": 912,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Дискование 2-е",
        "culture": "Свекла сахарная",
        "per_day": 59,
        "per_operation": 1041,
        "val_day": "",
        "val_beginning": ""
      }
    ])
  },
  {
    "input": "12.10\nВнесение мин удобрений под оз пшеницу 2025 г ПУ Юг 149/7264\nОтд 17 -149/1443",
    "output": safe_json([
      {
        "date": "12-окт.",
        "division": "АОР",
        "operation": "Внесение минеральных удобрений",
        "culture": "Пшеница озимая товарная",
        "per_day": 149,
        "per_operation": 7264,
        "val_day": "",
        "val_beginning": ""
      }
    ])
  },
  {
    "input": "Север\nОтд7 пах с св 41/501\nОтд20 20/281 по пу 61/793\nОтд 3 пах подс.60/231\nПо пу 231\n\nДиск к. Сил отд 7. 32/352\nПу- 484\n\nДиск под Оз п езубов 20/281\nДиск под с. Св отд 10 83/203 пу-1065га",
    "output": safe_json([
      {
        "date": "",
        "division": "АОР",
        "operation": "Пахота",
        "culture": "Свекла сахарная",
        "per_day": 61,
        "per_operation": 793,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Пахота",
        "culture": "Подсолнечник товарный",
        "per_day": 60,
        "per_operation": 231,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Дискование",
        "culture": "Кукуруза кормовая",
        "per_day": 32,
        "per_operation": 484,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Дискование",
        "culture": "Пшеница озимая товарная",
        "per_day": 20,
        "per_operation": 281,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Дискование",
        "culture": "Свекла сахарная",
        "per_day": 83,
        "per_operation": 203,
        "val_day": "",
        "val_beginning": ""
      }
    ])
  },
  {
    "input": "Внесение удобрений под рапс отд 7\n-138/270\nДискование под рапс 40/172\nДиск после Кук сил отд 7 - 32/352 по пу 484га",
    "output": safe_json([
      {
        "date": "",
        "division": "АОР",
        "operation": "Внесение минеральных удобрений",
        "culture": "Рапс озимый",
        "per_day": 138,
        "per_operation": 270,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Дискование",
        "culture": "Рапс озимый",
        "per_day": 40,
        "per_operation": 172,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Дискование",
        "culture": "Кукуруза кормовая",
        "per_day": 32,
        "per_operation": 484,
        "val_day": "",
        "val_beginning": ""
      }
    ])
  },
  {
    "input": "10.03 день\n2-я подкормка озимых, ПУ \"Юг\" - 1749/2559\n(в т.ч Амазон - 1082/1371\nПневмоход - 667/1188)\n\nОтд11 - 307/307 (амазон 307/307)\nОтд 12 - 671/671 (амазон 318/318; пневмоход 353/353)\nОтд 16 - 462/1272 (амазон 148/437; пневмоход 314/835)\nОтд 17 - 309/309 (амазон 309/309)",
    "output": safe_json([
      {
        "date": "3/10/2024",
        "division": "АОР",
        "operation": "2-я подкормка",
        "culture": "Пшеница озимая товарная",
        "per_day": 1749,
        "per_operation": 2559,
        "val_day": "",
        "val_beginning": ""
      }
    ])
  },
  {
    "input": "Уборка свеклы 27.10. день\nОтд10 - 45/216\nПо ПУ 45/1569\nВал 1259680/6660630\nУрожайность 279.9/308.3\nПо ПУ 1259680/41630060\nНа завод 1811630/6430580\nПо ПУ 1811630/41400550\nПоложено в кагат 399400\nВвезено с кагата 951340\nОстаток 230060\nОз-9,04/12,58\nДигестия -14,50/15,05",
    "output": safe_json([
      {
        "date": "10/27/2024",
        "division": "АОР",
        "operation": "Уборка",
        "culture": "Свекла сахарная",
        "per_day": 45,
        "per_operation": 1569,
        "val_day": "12 596,80",
        "val_beginning": "66 606,30"
      }
    ])
  },
  {
    "input": "Пахота под сах св\nПо Пу 77/518\nОтд 12 46/298\nОтд 16 21/143\nОтд 17 10/17\n\nЧизел под оз ячмень\nПо Пу 22/640\nОтд 11 22/242\n\nЧизел под оз зел корм\nОтд 11 40/40\n\nДиск оз пшеницы\nПо Пу 28/8872\nОтд 17 28/2097\n\n2-е диск под сах св\nПо Пу 189/1763\nОтд 11 60/209\nОтд 12 122/540\nОтд 17 7/172\n\nДиск кук силос\nПо Пу 6/904\nОтд 11 6/229\n\nПрик под оз ячмень\nПо Пу 40/498\nОтд 11 40/100\n\nУборка сои (семенной)\nОтд 11 65/65\nВал 58720\nУрож 9",
    "output": safe_json([
      {
        "date": "",
        "division": "Юг",
        "operation": "Пахота",
        "culture": "Свекла сахарная",
        "per_day": 77,
        "per_operation": 518,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Чизелевание",
        "culture": "Ячмень озимый",
        "per_day": 22,
        "per_operation": 640,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Чизелевание",
        "culture": "Пшеница озимая на зеленку",
        "per_day": 40,
        "per_operation": 40,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Дискование",
        "culture": "Пшеница озимая ",
        "per_day": 28,
        "per_operation": 8872,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Дискование 2-е",
        "culture": "Свекла сахарная",
        "per_day": 189,
        "per_operation": 1763,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Дискование",
        "culture": "Кукуруза кормовая",
        "per_day": 6,
        "per_operation": 904,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Прикатывание посевов",
        "culture": "Ячмень озимый",
        "per_day": 40,
        "per_operation": 498,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Уборка",
        "culture": "Соя семенная",
        "per_day": 65,
        "per_operation": 65,
        "val_day": "587,20",
        "val_beginning": ""
      }
    ])
  },
  {
    "input": "Пахота под сах св\nПоПу 91/609\nОтд 11 13/73\nОтд 12 49/347\nОтд 16 20/163\nОтд 17 9/26\n\nЧизел под оз зел корм\nОтд 11 60/100\n\n2-е диск под сах св\nПо Пу 53/1816\nОтд 12 53/593\n\nДиск кук силос\nПо Пу 66/970\nОтд 11 66/295\n\nДиск сах св\nОтд 12 13/13\n\nПрикат под оз ячмень\nПо Пу 40/538\nОтд 11 40/140\n\nУборка сои семенной\nОтд 11 29/94\nВал 37400/96120\nУрож 12,9/10,2",
    "output": safe_json([
      {
        "date": "",
        "division": "Юг",
        "operation": "Пахота",
        "culture": "Свекла сахарная",
        "per_day": 91,
        "per_operation": 609,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Чизелевание",
        "culture": "Пшеница озимая на зеленку",
        "per_day": 60,
        "per_operation": 100,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Дискование 2-е",
        "culture": "Свекла сахарная",
        "per_day": 53,
        "per_operation": 1816,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Дискование",
        "culture": "Кукуруза кормовая",
        "per_day": 66,
        "per_operation": 970,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Дискование",
        "culture": "Свекла сахарная",
        "per_day": 13,
        "per_operation": 13,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Прикатывание посевов",
        "culture": "Ячмень озимый",
        "per_day": 40,
        "per_operation": 538,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Уборка",
        "culture": "Соя семенная",
        "per_day": 29,
        "per_operation": 94,
        "val_day": "374,00",
        "val_beginning": "961,20"
      }
    ])
  },
  {
    "input": "Пахота под сах св\nПо Пу 88/329\nОтд 11 33/60\nОтд 12 34/204\nОтд 16 21/65\n\nПахота под мн тр\nПо Пу 10/438\nОтд 17 10/80\n\nЧизел под оз ячмень\nПо Пу 71/528\nОтд 11 71/130\n\n2-е диск под сах св\nПо Пу 80/1263\nОтд 12 80/314\n\n2-е диск под оз ячмень\nПо Пу 97/819\nОтд 17 97/179\n\nДиск кук силос\nПо Пу 43/650\nОтд 11 33/133\nОтд 12 10/148\n\nВышка отц форм под/г\nОтд 12 10/22\n\nУборка сах св\nОтд 12 16/16\nВал 473920\nУрож 296.2\nДиг - 19,19\nОз - 5,33",
    "output": safe_json([
      {
        "date": "",
        "division": "АОР",
        "operation": "Пахота",
        "culture": "Свекла сахарная",
        "per_day": 88,
        "per_operation": 329,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Пахота",
        "culture": "Многолетние травы текущие",
        "per_day": 10,
        "per_operation": 438,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Чизелевание",
        "culture": "Ячмень озимый",
        "per_day": 71,
        "per_operation": 528,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Дискование 2-е",
        "culture": "Свекла сахарная",
        "per_day": 80,
        "per_operation": 1263,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Дискование 2-е",
        "culture": "Ячмень озимый",
        "per_day": 97,
        "per_operation": 819,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Дискование",
        "culture": "Кукуруза кормовая",
        "per_day": 43,
        "per_operation": 650,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "АОР",
        "operation": "Уборка",
        "culture": "Свекла сахарная",
        "per_day": 16,
        "per_operation": 16,
        "val_day": "4 739,20",
        "val_beginning": ""
      }
    ])
  },
  {
    "input": "20.11 Мир\nПахота зяби под сою 100 га день, 1109 га от начала, 97%, 30 га остаток.\nРаботало 4 агрегата.\n\nВыравнивание зяби под подсолнечник\n47 га день, 141 га от начала, 29 %, остаток 565 га. Работал 1 агрегат\n\nОсадки:\nБригада 1 Воронежская – 6 мм",
    "output": safe_json([
      {
        "date": "20/11",
        "division": "Мир",
        "operation": "Пахота",
        "culture": "Соя товарная",
        "per_day": 100,
        "per_operation": 1109,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "21/11",
        "division": "Мир",
        "operation": "Выравнивание зяби",
        "culture": "Подсолнечник товарный",
        "per_day": 47,
        "per_operation": 141,
        "val_day": "",
        "val_beginning": ""
      }
    ])
  },
  {
    "input": "ТСК\nВыравнивание зяби под сою 25 га/ с нарастающим 765 га (13%) Остаток 5332 га\nВыравнивание зяби под кукурузу 131 га (3%) Остаток 4486 га\nОсадки 1мм",
    "output": safe_json([
      {
        "date": "",
        "division": "ТСК",
        "operation": "Пахота",
        "culture": "Соя товарная",
        "per_day": 25,
        "per_operation": 765,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "ТСК",
        "operation": "Выравнивание зяби",
        "culture": "Кукуруза товарная",
        "per_day": 131,
        "per_operation": 131,
        "val_day": "",
        "val_beginning": ""
      }
    ])
  },
  {
    "input": "Восход\nПосев кук - 24/252 га\n24%\nПредпосевная культ под кук - 94/490 га 46%\nСЗР оз пш - 103/557 га\n25%\nПодкормка оз рапс - 152 га, 100%, подкормка овса - 97 га, 50%\nДовсходовое боронование подсолнечника - 524 га, 100%",
    "output": safe_json([
      {
        "date": "",
        "division": "Восход",
        "operation": "Сев",
        "culture": "Кукуруза товарная",
        "per_day": 24,
        "per_operation": 252,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "Восход",
        "operation": "Предпосевная культивация",
        "culture": "Кукуруза товарная",
        "per_day": 94,
        "per_operation": 490,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "Восход",
        "operation": "Гербицидная обработка",
        "culture": "Пшеница озимая товарная",
        "per_day": 103,
        "per_operation": 557,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "Восход",
        "operation": "Подкормка",
        "culture": "Рапс озимый",
        "per_day": 152,
        "per_operation": 152,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "Восход",
        "operation": "Подкормка",
        "culture": "Овес",
        "per_day": 97,
        "per_operation": 97,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "",
        "division": "Восход",
        "operation": "Боронование довсходовое",
        "culture": "Подсолнечник товарный",
        "per_day": 524,
        "per_operation": 524,
        "val_day": "",
        "val_beginning": ""
      }
    ])
  },
  {
    "input": "30.03.25 г.\nСП Коломейцево\nпредпосевная культивация под подсолнечник — день 30 га, от начала 187 га (91%)\nсев подсолнечника — день+ночь 57 га, от начала 157 га (77%)\nВнесение почвенного гербицида по подсолнечнику — день 82 га, от начала 82 га (38%)",
    "output": safe_json([
      {
        "date": "30/3",
        "division": "СП Коломейцево",
        "operation": "Предпосевная культивация",
        "culture": "Подсолнечник товарный",
        "per_day": 30,
        "per_operation": 187,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "30/3",
        "division": "СП Коломейцево",
        "operation": "Сев",
        "culture": "Подсолнечник товарный",
        "per_day": 57,
        "per_operation": 157,
        "val_day": "",
        "val_beginning": ""
      },
      {
        "date": "30/3",
        "division": "СП Коломейцево",
        "operation": "Гербицидная обработка",
        "culture": "Подсолнечник товарный",
        "per_day": 82,
        "per_operation": 82,
        "val_day": "",
        "val_beginning": ""
      }
    ])
  }
]  