# Parenthesis Service API 

## Оглавление:

**[1. Запуск проекта локально](#local)**  
**[2. API](#api)**  
&emsp;**[2.1 Validate](#validate)**  
&emsp;**[2.2 Fix](#fix)**

### 1. Запуск проекта локально

```make run``` - запускает сервис на порту 50051, прометеус на порту 9090

```make down``` - останавливает и удаляет контейнеры    


<a name="api"></a>

### 2. API

Два метода API, ожидаю на вход строчную последовательность из скобок.
1) Validate – проверяет правильны ли пары и порядки
2) Fix – закрывает скобки или удаляет лишние

<a name="validate"></a>

### 2.1 Validate

**REQUEST EXAMPLE:**   
```text
  [()]{}{[()()]()}
```

**RESPONSE EXAMPLE:**   
```json
{
  "result": "Balanced"
}
```

<a name="fix"></a>

### 2.2 Fix

**REQUEST EXAMPLE:**   
```text
  ()[]{((())){}{[]}
```

**RESPONSE EXAMPLE:**   
```json
{
  "strOut": "()[]{((())){}{[]}}"
}
```

