package config

var (
	Cfg Config
)

type Config struct {
	SecondsSaveStats uint // сколько секунд храним статистику в памяти
	ClearStatsSecondsInterval uint // раз во сколько секунд очищаем статистику от старых данных
}

func Load() {
	Cfg = Config{
		SecondsSaveStats: 300,
	}
}
