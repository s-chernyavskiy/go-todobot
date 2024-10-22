package postgresql

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-todobot/internal/storage"
	"log"
)

type Postgresql struct {
	connPool *pgxpool.Pool
}

func New(conn *pgxpool.Pool) *Postgresql {
	return &Postgresql{
		connPool: conn,
	}
}

func (p Postgresql) AddTask(t *storage.Task) error {
	rows, err := p.connPool.Query(context.Background(),
		`insert into tasks (user_id, name) values ($1, $2)`,
		t.UserID,
		t.Name,
	)

	defer rows.Close()

	return err
}

func (p Postgresql) RemoveTask(t *storage.Task) error {
	log.Println("started to remove task from db")
	rows, err := p.connPool.Query(context.Background(),
		`delete from tasks where (name=$1 and user_id=$2)`, t.Name, t.UserID,
	)

	defer rows.Close()

	log.Println("finished removing task from db")
	return err
}

func (p Postgresql) ListTasks(userId string) ([]string, error) {
	r, err := p.connPool.Query(context.Background(), "select * from tasks where user_id=$1", userId)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	var tasks []string

	for r.Next() {
		var t storage.Task

		err := r.Scan(&t.TaskID, &t.UserID, &t.Name)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t.Name)
	}

	return tasks, r.Err()
}

func (p Postgresql) IfExistsTask(t *storage.Task) (bool, error) {
	r, err := p.connPool.Query(context.Background(), "select * from tasks where name=$1", t.Name)
	if err != nil {
		return false, err
	}

	defer r.Close()

	return r.Next(), nil
}
