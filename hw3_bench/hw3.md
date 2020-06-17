Есть функиця, которая что-то там ищет по файлу. Но делает она это не очень быстро. Надо её оптимизировать.

Задание на работу с профайлером pprof.

Цель задания - научиться работать с pprof, находить горячие места в коде, уметь строить профиль потребления cpu и памяти, оптимизировать код с учетом этой информации. Написание самого быстрого решения не является целью задания.

Для генерации графа вам понадобится graphviz. Для пользователей windows не забудьте добавить его в PATH чтобы была доступна команда dot.

Рекомендую внимательно прочитать доп. материалы на русском - там ещё много примеров оптимизации и объяснений как работать с профайлером. Фактически там есть вся информация для выполнения этого задания.

Есть с десяток мест где можно оптимизировать.

Для выполнения задания необходимо чтобы один из параметров ( ns/op, B/op, allocs/op ) был быстрее чем в *BenchmarkSolution* ( fast < solution ) и ещё один лучше *BenchmarkSolution* + 20% ( fast < solution * 1.2), например ( fast allocs/op < 10422*1.2=12506 ).

По памяти ( B/op ) и количеству аллокаций ( allocs/op ) можно ориентироваться ровно на результаты *BenchmarkSolution* ниже, по времени ( ns/op ) - нет, зависит от системы.

Для этого задания увеличено количество проверок с 3-х до 5 за 8 часов.

Результат в fast.go в функцию FastSearch (изначально там то же самое что в SlowSearch).

Пример результатов с которыми будет сравниваться:
```
$ go test -bench . -benchmem

goos: windows

goarch: amd64

BenchmarkSlow-8 10 142703250 ns/op 336887900 B/op 284175 allocs/op

BenchmarkSolution-8 500 2782432 ns/op 559910 B/op 10422 allocs/op

PASS

ok coursera/hw3 3.897s
```

Запуск:
* `go test -v` - чтобы проверить что ничего не сломалось
* `go test -bench . -benchmem` - для просмотра производительности

Советы:
* Смотрите где мы аллоцируем память
* Смотрите где мы накапливаем весь результат, хотя нам все значения одновременно не нужны
* Смотрите где происходят преобразования типов, которые можно избежать
* Смотрите не только на графе, но и в pprof в текстовом виде (list FastSearch) - там прямо по исходнику можно увидеть где что
* Задание предполагает использование easyjson. На сервере эта библиотека есть, подключать можно. Но сгенерированный через easyjson код вам надо поместить в файл с вашей функцией
* Можно сделать без easyjson

Примечание:
* easyjson основан на рефлекции и не может работать с пакетом main. Для генерации кода вам необходимо вынести вашу структуру в отдельный пакет, сгенерить там код, потом забрать его в main


###Результат
```
[sim@oxy hw3_bench] (master) $ make
go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1 .
goos: darwin
goarch: amd64
pkg: coursera/hw3_bench
BenchmarkSlow-4   	       2	 875346638 ns/op	18808992 B/op	  195812 allocs/op
BenchmarkFast-4   	      32	  34618589 ns/op	  707043 B/op	   10197 allocs/op
PASS
ok  	coursera/hw3_bench	4.802s
go tool pprof hw3_bench.test cpu.out
File: hw3_bench.test
Type: cpu
Time: Jun 17, 2020 at 2:40pm (EEST)
Duration: 4.69s, Total samples = 4.54s (96.79%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 3630ms, 79.96% of 4540ms total
Dropped 92 nodes (cum <= 22.70ms)
Showing top 10 nodes out of 114
      flat  flat%   sum%        cum   cum%
     880ms 19.38% 19.38%      930ms 20.48%  runtime.addspecial
     820ms 18.06% 37.44%      820ms 18.06%  syscall.syscall
     520ms 11.45% 48.90%      620ms 13.66%  runtime.step
     380ms  8.37% 57.27%     1020ms 22.47%  runtime.pcvalue
     320ms  7.05% 64.32%      320ms  7.05%  runtime.madvise
     210ms  4.63% 68.94%     1480ms 32.60%  runtime.gentraceback
     200ms  4.41% 73.35%     1180ms 25.99%  runtime.setprofilebucket
     110ms  2.42% 75.77%      130ms  2.86%  runtime.findfunc
     100ms  2.20% 77.97%      100ms  2.20%  runtime.readvarint
      90ms  1.98% 79.96%       90ms  1.98%  runtime.usleep
(pprof) web
```
