TODO: конфиг, чтобы отсекать юнит-тесты

Давайте посмотрим на код анализатора и возможные оптимизации:
Read file: main.go
Вижу несколько возможных оптимизаций:

1. **Параллельная обработка файлов**:
```go
func collectInterfaceMethods(pass *analysis.Pass) map[*types.Func]methodInfo {
    ifaceMethods := sync.Map{}
    var wg sync.WaitGroup
    
    for _, file := range pass.Files {
        wg.Add(1)
        go func(file *ast.File) {
            defer wg.Done()
            // Анализ файла
            // Сохранение результатов в ifaceMethods
        }(file)
    }
    wg.Wait()
    
    // Конвертация в обычную map
    result := make(map[*types.Func]methodInfo)
    ifaceMethods.Range(func(key, value interface{}) bool {
        result[key.(*types.Func)] = value.(methodInfo)
        return true
    })
    return result
}
```

Результат и анализ
Параллельная обработка файлов реализована, но эффекта почти нет. Вот что показал профилер:
Где тратится время:
49.63% — системные вызовы (runtime.cgocall)
28-30% — I/O операции (чтение файлов)
15% — парсинг и проверка типов Go
<1% — наша функция collectInterfaceMethods
Почему параллелизация не сработала:
Функция не узкое место — collectInterfaceMethods занимает ничтожную долю времени
Мало файлов — в тестах всего 3 файла, overhead горутин больше выигрыша
I/O доминирует — основное время уходит на чтение/парсинг файлов
Реальные узкие места:
Чтобы ускорить анализатор нужно:
Кешировать результаты парсинга — 15% времени
Параллелить I/O операции — 30% времени
Оптимизировать системные вызовы — 50% времени



2. **Оптимизация поиска методов**:
```go
// Предварительное кэширование имен методов
type methodAnalyzer struct {
    // ... существующие поля ...
    methodNames map[string][]*types.Func // кэш методов по именам
}

func (ma *methodAnalyzer) buildMethodNameCache() {
    ma.methodNames = make(map[string][]*types.Func)
    for method := range ma.ifaceMethods {
        name := method.Name()
        ma.methodNames[name] = append(ma.methodNames[name], method)
    }
}
```

3. **Раннее прерывание проверок**:
```go
func (ma *methodAnalyzer) markMatchingMethods(calledMethod *types.Func, recv types.Type) {
    // Сначала проверяем по имени
    candidates := ma.methodNames[calledMethod.Name()]
    if len(candidates) == 0 {
        return // быстрый выход если нет методов с таким именем
    }
    
    // Дальше проверяем только подходящие кандидаты
    for _, ifaceMethod := range candidates {
        if ma.usedMethods[ifaceMethod] {
            continue
        }
        // ... остальные проверки
    }
}
```

4. **Оптимизация проверки реализации интерфейса**:
```go
type methodAnalyzer struct {
    // ... существующие поля ...
    implementsCache map[implementsCacheKey]bool
}

type implementsCacheKey struct {
    recv     types.Type
    iface    *types.Interface
    hashCode uint64 // для быстрого сравнения
}

func (ma *methodAnalyzer) checkImplements(recv types.Type, info methodInfo) bool {
    key := implementsCacheKey{recv: recv, iface: info.iface}
    if result, ok := ma.implementsCache[key]; ok {
        return result
    }
    
    result := types.Implements(recv, info.iface)
    ma.implementsCache[key] = result
    return result
}
```

5. **Уменьшение аллокаций**:
```go
// Переиспользование слайсов и мап
var nodeFilter = []ast.Node{
    (*ast.SelectorExpr)(nil),
    (*ast.CallExpr)(nil),
}

func (ma *methodAnalyzer) analyze() map[*types.Func]bool {
    ins := ma.pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
    ins.Preorder(nodeFilter, ma.analyzeNode)
    return ma.usedMethods
}
```

Какую из этих оптимизаций хотите попробовать первой? Или есть другие идеи?