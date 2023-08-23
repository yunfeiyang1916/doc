package service

import (
	"context"
	"log"
	"testing"
	"user_growth/conf"
	"user_growth/dbhelper"
	"user_growth/models"

	_ "github.com/go-sql-driver/mysql"
)

func initDb() {
	// 统一设置时区为UTC 零时区
	//time.Local = time.UTC
	conf.LoadConfigs()
	dbhelper.InitDb()
}

func TestCoinTaskService_Save(t *testing.T) {
	initDb()
	s := NewCoinTaskService(context.Background())
	data := models.TbCoinTask{
		Id:    0,
		Task:  "post article",
		Coin:  10,
		Limit: 10,
	}
	if err := s.Save(&data); err != nil {
		t.Errorf("Save(%+v) error=%v", data, err)
	} else {
		log.Printf("Save data=%+v\n", data)
	}
}

func TestCoinTaskService_GetByTask(t *testing.T) {
	initDb()
	s := NewCoinTaskService(context.Background())
	task := "post article"
	if data, err := s.GetByTask(task); err != nil {
		t.Errorf("GetByTask(%s) error=%v", task, err)
	} else {
		log.Printf("GetByTask(%s) data=%+v", task, data)
	}
}

func TestCoinTaskService_FindAll(t *testing.T) {
	initDb()
	s := NewCoinTaskService(context.Background())
	if datalist, err := s.FindAll(); err != nil {
		t.Errorf("FindAll() error=%v", err)
	} else {
		log.Printf("FindAll data=%+v", datalist)
	}
}

func TestGradeInfoService_Save(t *testing.T) {
	initDb()
	s := NewGradeInfoService(context.Background())
	data := models.TbGradeInfo{
		Id:          0,
		Title:       "初级",
		Description: "初级用户",
		Score:       0,
		Expired:     0,
	}
	if err := s.Save(&data); err != nil {
		t.Errorf("Save(%+v) error=%v", data, err)
	} else {
		log.Print("看看换行了没")
		log.Printf("Save(%+v) success", data)
	}
}

func TestGradeInfoService_Get(t *testing.T) {
	initDb()
	s := NewGradeInfoService(context.Background())
	if data, err := s.Get(1); err != nil {
		t.Errorf("Get(1) error=%v", err)
	} else {
		log.Printf("Get(1) data=%+v\n", data)
	}
}

func TestGradeInfoService_FindAll(t *testing.T) {
	initDb()
	s := NewGradeInfoService(context.Background())
	if datalist, err := s.FindAll(); err != nil {
		t.Errorf("FindAll error=%v", err)
	} else {
		log.Printf("FindAll datalist=%+v\n", datalist)
	}
}

func TestGradeInfoService_NowGrade(t *testing.T) {
	initDb()
	s := NewGradeInfoService(context.Background())
	if data, err := s.NowGrade(1); err != nil {
		t.Errorf("NowGrade error=%v", err)
	} else {
		log.Printf("NowGrade v=%+v\n", data)
	}
}
