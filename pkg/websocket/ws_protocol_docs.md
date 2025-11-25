## WebSocket Message Protocol

### Структура сообщения

Все WebSocket сообщения имеют единую структуру:

```json
{
  "type": "string",           // Тип сообщения (см. ниже)
  "payload": {},              // Данные сообщения (зависят от типа)
  "timestamp": "ISO8601",     // Время отправки
  "agent_id": "string",       // ID агента (опционально)
  "request_id": "string"      // ID запроса для корреляции (опционально)
}
```

### Типы сообщений

#### 1. `heartbeat` (Agent → Backend)

**Назначение:** Периодическая отправка статуса агента (каждые 30 секунд)

**Payload:**
```json
{
  "status": "online|idle|busy",
  "cpu_usage": 23.5,
  "memory_usage": 45.8,
  "username": "current_user"
}
```

**Маршрутизация:** Backend транслирует в канал `admin::dashboard` как `agent_status`

---

#### 2. `status_update` (Agent → Backend)

**Назначение:** Уведомление об изменении статуса агента

**Payload:**
```json
{
  "old_status": "idle",
  "new_status": "busy",
  "reason": "User logged in"
}
```

**Маршрутизация:** Broadcast в `admin::dashboard`

---

#### 3. `action_log` (Agent → Backend)

**Назначение:** Логирование действий пользователя (блокировки процессов, URL, доступ к файлам)

**Payload:**
```json
{
  "action": "process_blocked|url_blocked|file_accessed",
  "resource": "chrome.exe|example.com|C:\\file.txt",
  "username": "student01",
  "allowed": false,
  "metadata": {
    "rule_id": 123,
    "group_dn": "CN=Students,..."
  }
}
```

**Маршрутизация:**
- Сохранение в БД (`agent_logs`)
- Broadcast в `admin::dashboard`

---

#### 4. `command_request` (Admin → Backend → Agent)

**Назначение:** Отправка команды от админа к конкретному агенту

**Payload:**
```json
{
  "command": "get_processes|kill_process|update_whitelist",
  "params": {
    "process_name": "chrome.exe",
    "pid": 1234
  },
  "timeout": 30
}
```

**Маршрутизация:** Backend направляет в канал `agent::{agent_id}`

**Требования:**
- Только пользователи с ролью `admin` могут отправлять команды
- `request_id` обязателен для корреляции с ответом

---

#### 5. `command_response` (Agent → Backend → Admin)

**Назначение:** Ответ агента на команду от админа

**Payload:**
```json
{
  "success": true,
  "result": {
    "processes": [...]
  },
  "error": "Optional error message"
}
```

**Маршрутизация:** Backend возвращает в канал `admin::dashboard` (TODO: таргетинг к конкретному админу)

**Требования:**
- `request_id` должен совпадать с исходной командой

---

#### 6. `agent_status` (Backend → Admin)

**Назначение:** Уведомление админов об изменении статуса агентов

**Payload:**
```json
{
  "agent_id": "agent-001",
  "hostname": "PC-LAB-01",
  "status": "online|offline",
  "last_seen": "ISO8601",
  "username": "student01"
}
```

**Генерируется:** Backend при получении `heartbeat` от агента

---

### Каналы WebSocket

| Канал | Участники | Назначение |
|-------|-----------|------------|
| `agent::{agent_id}` | Конкретный агент | Прием команд от админов |
| `admin::dashboard` | Все админы | Мониторинг всех агентов, получение логов |

### Reconnection Logic

**Exponential Backoff:**
- Начальная задержка: 1 секунда
- Максимальная задержка: 30 секунд
- Формула: `delay = min(baseDelay * 2^attempts, maxDelay)`

**Примеры задержек:**
- Попытка 1: 1 сек
- Попытка 2: 2 сек
- Попытка 3: 4 сек
- Попытка 4: 8 сек
- Попытка 5: 16 сек
- Попытка 6+: 30 сек

**Сброс счетчика:** При успешном подключении `reconnectAttempts = 0`

---

### Примеры использования

#### Агент отправляет heartbeat

```go
payload := HeartbeatPayload{
    Status:      "online",
    CPUUsage:    15.2,
    MemoryUsage: 42.8,
    Username:    "student01",
}
client.Send(MessageTypeHeartbeat, payload, "agent-001")
```

#### Админ отправляет команду агенту

```javascript
const payload = {
    command: "get_processes",
    params: {},
    timeout: 30
};
client.send("command_request", payload, "agent-001");
```

#### Агент обрабатывает команду

```go
client.On(MessageTypeCommandRequest, func(msg *Message) error {
    var payload CommandRequestPayload
    msg.DecodePayload(&payload)
    
    // Execute command
    result := executeCommand(payload.Command, payload.Params)
    
    // Send response
    response := CommandResponsePayload{
        Success: true,
        Result:  result,
    }
    return client.Send(MessageTypeCommandResponse, response, msg.AgentID)
})
```