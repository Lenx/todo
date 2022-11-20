package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/lenx/todo/pkg/handler"
	"github.com/lenx/todo/pkg/repository"
	"github.com/lenx/todo/pkg/service"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func main() {

	var c conf
	c.getConf()

	//fmt.Println(c.DB.Port)

	// задаём формат логов
	logrus.SetFormatter(new(logrus.JSONFormatter))
	/*
		if err := initConfig(); err != nil {
			logrus.Fatal("error initializing configs:", err.Error())
		}
	*/
	// подгружаем .env файл
	if err := godotenv.Load(); err != nil {
		logrus.Fatal("error loading env variable: ", err.Error())
	}

	// запускаем базу данных
	/*
		db, err := repository.NewPostrgesDB(repository.Config{
			Host: viper.GetString("db.host"),
			Port: viper.GetString("db.port"),
			Username: viper.GetString("db.username"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName: viper.GetString("db.dbname"),
			SSLMode: viper.GetString("db.sslmode"),
		})
	*/
	db, err := repository.NewPostrgesDB(repository.Config{
		Host:     c.DB.Host,
		Port:     c.DB.Port,
		Username: c.DB.Username,
		Password: c.DB.Password,
		DBName:   c.DB.Dbname,
		//SSLMode: c.DB.Sslmode,
	})
	if err != nil {
		logrus.Fatal("failed to initialize db: ", err.Error())

	}

	// устанавливаем зависимости
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	// запускаем сервер
	srv := new(todo_app.Server)
	go func() {
		if err := srv.Run(os.Getenv("PORT"), handlers.InitRoutes()); err != nil {
			logrus.Fatal("error with running http server:", err.Error())

		}
	}()

	logrus.Print("TodoApp Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("TodoApp shootting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occurred on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occurred on db connection close: %s", err.Error())
	}

}

// читаем config
/*
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
*/

type conf struct {
	Port string `yaml:"port"`
	DB   struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Dbname   string `yaml:"dbname"`
		Sslmode  string `yaml:"sslmode"`
	}
}

func (c *conf) getConf() *conf {

	yamlFile, err := ioutil.ReadFile("configs/config.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
