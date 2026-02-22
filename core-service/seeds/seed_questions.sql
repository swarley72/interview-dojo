-- Seed: 100 вопросов с разными type, difficulty, tags
-- Запуск: psql -f seeds/seed_questions.sql "$DATABASE_URL"

BEGIN;

-- 1. Теги
INSERT INTO tags (name) VALUES
    ('go'), ('concurrency'), ('goroutines'), ('channels'), ('context'),
    ('interfaces'), ('generics'), ('testing'), ('benchmarks'), ('stdlib'),
    ('sql'), ('postgresql'), ('redis'), ('kafka'), ('grpc'),
    ('protobuf'), ('docker'), ('kubernetes'), ('ci-cd'), ('git'),
    ('algorithms'), ('data-structures'), ('trees'), ('graphs'), ('sorting'),
    ('system-design'), ('microservices'), ('caching'), ('load-balancing'), ('api-design'),
    ('networking'), ('http'), ('tcp'), ('dns'), ('tls'),
    ('linux'), ('os'), ('memory'), ('garbage-collector'), ('scheduler')
ON CONFLICT (name) DO NOTHING;

-- 2. Вопросы (100 шт.)
-- Генерируем через generate_series + массивы для рандомных комбинаций

INSERT INTO questions (id, type, title, content_md, answer_md, difficulty)
SELECT
    gen_random_uuid(),
    -- type: один из 4
    (ARRAY['theory', 'coding', 'system_design', 'algorithm'])[1 + (i % 4)],
    -- title
    CASE (i % 4)
        WHEN 0 THEN (ARRAY[
            'Как работают горутины в Go?',
            'Объясните принцип работы каналов',
            'Что такое context.Context и зачем он нужен?',
            'Разница между mutex и rwmutex',
            'Что такое интерфейсы в Go?',
            'Как устроен garbage collector в Go?',
            'Объясните модель памяти Go',
            'Что такое defer и как он работает?',
            'Как работает select в Go?',
            'Что такое race condition?',
            'Зачем нужен sync.WaitGroup?',
            'Как работает sync.Pool?',
            'Что такое embedding в Go?',
            'Как устроен scheduler в Go?',
            'Разница между слайсом и массивом',
            'Что такое zero value в Go?',
            'Как работает type assertion?',
            'Что такое stringer interface?',
            'Как устроен map в Go?',
            'Что такое panic и recover?',
            'Зачем нужен init() в Go?',
            'Как работает сборка мусора?',
            'Что такое escape analysis?',
            'Как профилировать Go-приложение?',
            'Что такое generics в Go?'
        ])[1 + (i / 4) % 25]
        WHEN 1 THEN (ARRAY[
            'Реализуйте worker pool на горутинах',
            'Напишите concurrent-safe кэш',
            'Реализуйте rate limiter',
            'Напишите graceful shutdown сервера',
            'Реализуйте pipeline из горутин',
            'Напишите middleware для HTTP-сервера',
            'Реализуйте retry с exponential backoff',
            'Напишите парсер JSON без encoding/json',
            'Реализуйте LRU-кэш',
            'Напишите простой HTTP-роутер',
            'Реализуйте fan-out/fan-in паттерн',
            'Напишите connection pool для БД',
            'Реализуйте circuit breaker',
            'Напишите unit-тесты с моками',
            'Реализуйте pub/sub на каналах',
            'Напишите CLI-утилиту с cobra',
            'Реализуйте binary search generics',
            'Напишите обёртку над pgx',
            'Реализуйте middleware chain',
            'Напишите бенчмарк для функции',
            'Реализуйте semaphore на каналах',
            'Напишите gRPC-сервер и клиент',
            'Реализуйте in-memory очередь',
            'Напишите health check endpoint',
            'Реализуйте сериализацию в protobuf'
        ])[1 + (i / 4) % 25]
        WHEN 2 THEN (ARRAY[
            'Спроектируйте систему уведомлений',
            'Спроектируйте URL shortener',
            'Спроектируйте чат-систему',
            'Спроектируйте систему авторизации',
            'Спроектируйте rate limiting сервис',
            'Спроектируйте CDN',
            'Спроектируйте task scheduler',
            'Спроектируйте event sourcing систему',
            'Спроектируйте API Gateway',
            'Спроектируйте систему логирования',
            'Спроектируйте distributed cache',
            'Спроектируйте message broker',
            'Спроектируйте систему мониторинга',
            'Спроектируйте service mesh',
            'Спроектируйте CI/CD пайплайн',
            'Спроектируйте feature flag систему',
            'Спроектируйте систему A/B тестов',
            'Спроектируйте поисковый движок',
            'Спроектируйте систему рекомендаций',
            'Спроектируйте distributed lock',
            'Спроектируйте load balancer',
            'Спроектируйте конфиг-сервер',
            'Спроектируйте систему деплоя',
            'Спроектируйте трейсинг систему',
            'Спроектируйте систему очередей'
        ])[1 + (i / 4) % 25]
        ELSE (ARRAY[
            'Реализуйте бинарный поиск',
            'Реализуйте быструю сортировку',
            'Найдите цикл в связном списке',
            'Реализуйте BFS для графа',
            'Реализуйте DFS для графа',
            'Найдите кратчайший путь (Dijkstra)',
            'Реализуйте стек на слайсах',
            'Реализуйте очередь на слайсах',
            'Найдите два числа с заданной суммой',
            'Реализуйте merge sort',
            'Проверьте сбалансированность скобок',
            'Реализуйте trie (префиксное дерево)',
            'Найдите максимальную подпоследовательность',
            'Реализуйте heap (кучу)',
            'Найдите медиану двух массивов',
            'Реализуйте hash map с нуля',
            'Найдите LCA в бинарном дереве',
            'Реализуйте топологическую сортировку',
            'Разверните связный список',
            'Найдите k-й элемент с конца списка',
            'Реализуйте union-find',
            'Сериализуйте бинарное дерево',
            'Найдите все перестановки строки',
            'Реализуйте sliding window',
            'Найдите longest common subsequence'
        ])[1 + (i / 4) % 25]
    END,
    -- content_md
    CASE (i % 4)
        WHEN 0 THEN format(
            E'## Теория\n\nОбъясните подробно следующую тему.\n\n### Требования\n\n- Дайте определение\n- Приведите примеры использования\n- Объясните внутреннее устройство\n- Укажите типичные ошибки\n\n> Уровень: %s',
            (ARRAY['junior', 'middle', 'senior'])[1 + (i % 3)]
        )
        WHEN 1 THEN format(
            E'## Задача\n\nНапишите решение на Go.\n\n### Ограничения\n\n- Используйте идиоматичный Go\n- Покройте edge cases\n- Напишите тесты\n\n```go\n// Ваш код здесь\nfunc Solution() {\n    // TODO\n}\n```\n\n> Сложность: %s',
            (ARRAY['O(n)', 'O(log n)', 'O(n log n)', 'O(1)'])[1 + (i % 4)]
        )
        WHEN 2 THEN format(
            E'## System Design\n\nСпроектируйте систему с учётом следующих требований.\n\n### Функциональные требования\n\n- Высокая доступность (99.9%%)\n- Масштабируемость до %s RPS\n- Latency < %s ms\n\n### Нефункциональные требования\n\n- Мониторинг и алертинг\n- Graceful degradation\n- Data consistency',
            (ARRAY['10K', '100K', '1M', '10M'])[1 + (i % 4)],
            (ARRAY['50', '100', '200', '500'])[1 + (i % 4)]
        )
        ELSE format(
            E'## Алгоритм\n\nРешите задачу оптимально.\n\n### Input\n\n```\n%s\n```\n\n### Output\n\nВерните результат в оптимальное время.\n\n### Constraints\n\n- 1 <= n <= %s\n- Элементы могут быть отрицательными',
            (ARRAY['nums = [2, 7, 11, 15], target = 9', 'head = [1, 2, 3, 4, 5]', 'graph = [[1,2],[2,3],[3,1]]', 's = "((()))"'])[1 + (i % 4)],
            (ARRAY['10^4', '10^5', '10^6', '10^9'])[1 + (i % 4)]
        )
    END,
    -- answer_md
    CASE (i % 4)
        WHEN 0 THEN format(
            E'## Ответ\n\n### Определение\n\nЭто ключевая концепция в Go, которая используется повсеместно в production-коде.\n\n### Как это работает\n\n1. Runtime создаёт структуру и управляет жизненным циклом\n2. Планировщик распределяет выполнение по OS-потокам\n3. GC отслеживает и освобождает неиспользуемую память\n\n### Пример\n\n```go\nfunc example() {\n    // Идиоматичный пример использования\n    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)\n    defer cancel()\n    \n    result, err := doWork(ctx)\n    if err != nil {\n        log.Error("failed", slog.String("error", err.Error()))\n        return\n    }\n    fmt.Println(result)\n}\n```\n\n### Частые ошибки\n\n- Забывать вызывать cancel() для context\n- Не обрабатывать ошибки\n- Игнорировать race conditions\n\n> **Собеседование:** часто спрашивают на уровне %s',
            (ARRAY['junior', 'middle', 'senior'])[1 + (i % 3)]
        )
        WHEN 1 THEN E'## Решение\n\n```go\npackage main\n\nimport (\n    "fmt"\n    "sync"\n)\n\nfunc Solution(input []int) []int {\n    mu := sync.Mutex{}\n    result := make([]int, 0, len(input))\n    \n    var wg sync.WaitGroup\n    for _, v := range input {\n        wg.Add(1)\n        go func(val int) {\n            defer wg.Done()\n            processed := val * 2\n            mu.Lock()\n            result = append(result, processed)\n            mu.Unlock()\n        }(v)\n    }\n    wg.Wait()\n    return result\n}\n```\n\n### Сложность\n\n- Время: O(n)\n- Память: O(n)\n\n### Тесты\n\n```go\nfunc TestSolution(t *testing.T) {\n    got := Solution([]int{1, 2, 3})\n    if len(got) != 3 {\n        t.Errorf("expected 3 elements, got %d", len(got))\n    }\n}\n```'
        WHEN 2 THEN E'## Архитектура\n\n### Компоненты\n\n```\nClient → Load Balancer → API Gateway → Service Mesh\n                                          ├── Service A (Go)\n                                          ├── Service B (Go)\n                                          └── Service C (Go)\n```\n\n### Хранилища\n\n- **PostgreSQL** — основные данные, ACID\n- **Redis** — кэш, сессии, rate limiting\n- **Kafka** — асинхронная коммуникация между сервисами\n- **S3** — файлы и бэкапы\n\n### Масштабирование\n\n1. Horizontal scaling через K8s\n2. Read replicas для PostgreSQL\n3. Redis Cluster для кэша\n4. Партиционирование Kafka топиков\n\n### Мониторинг\n\n- Prometheus + Grafana для метрик\n- Jaeger для distributed tracing\n- ELK для логов'
        ELSE E'## Решение\n\n```go\nfunc solve(nums []int, target int) []int {\n    seen := make(map[int]int)\n    for i, num := range nums {\n        if j, ok := seen[target-num]; ok {\n            return []int{j, i}\n        }\n        seen[num] = i\n    }\n    return nil\n}\n```\n\n### Разбор\n\n1. Используем hash map для O(1) lookup\n2. За один проход находим пару\n3. Обрабатываем edge case с пустым входом\n\n### Сложность\n\n- **Время:** O(n) — один проход по массиву\n- **Память:** O(n) — hash map\n\n### Альтернативы\n\n- Brute force: O(n²) время, O(1) память\n- Сортировка + два указателя: O(n log n) время, O(1) память'
    END,
    -- difficulty
    (ARRAY['easy', 'medium', 'hard'])[1 + (i % 3)]
FROM generate_series(1, 100) AS s(i);

-- 3. Привязка тегов к вопросам (2-4 тега на вопрос)
INSERT INTO question_tags (question_id, tag_id)
SELECT q.id, t.id
FROM (
    SELECT id, row_number() OVER (ORDER BY created_at, id) AS rn
    FROM questions
    ORDER BY created_at, id
) q
CROSS JOIN LATERAL (
    SELECT id FROM tags
    ORDER BY md5(q.id::text || tags.id::text)  -- детерминированный но "случайный" порядок
    LIMIT 2 + (q.rn % 3)  -- 2, 3 или 4 тега
) t
ON CONFLICT DO NOTHING;

COMMIT;
