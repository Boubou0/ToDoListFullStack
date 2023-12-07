package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", "./todolist.db")
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}

	DB = db
	return nil
}

type Task struct {
	Id          int    `json:"idTache"`
	IdTodoList  int    `json:"idTodoList"`
	Nom         string `json:"nom"`
	Description string `json:"description"`
	Etat        int    `json:"etat"`
	Priorite    int    `json:"priorite"`
	Deadline    string `json:"deadline"`
}

type Todolist struct {
	Id     int    `json:"idTodoLists"`
	Nom    string `json:"nom"`
	IdUser int    `json:"idUser"`
}

func GetTodoLists(userID uint) ([]Todolist, error) {
	rows, err := DB.Query("SELECT idTodoLists, nom from TODOLISTS WHERE idUser = ?", userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	liste := make([]Todolist, 0)

	for rows.Next() {
		singleTodolist := Todolist{}
		err = rows.Scan(&singleTodolist.Id, &singleTodolist.Nom)

		if err != nil {
			return nil, err
		}

		liste = append(liste, singleTodolist)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return liste, err
}

func GetTodoListById(id int) ([]Task, error) {

	rows, err := DB.Query("SELECT idTache, nom, description, etat, priorite, deadline from TACHE WHERE idTodoList =" + strconv.Itoa(id))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	listeTask := make([]Task, 0)

	for rows.Next() {
		singleTask := Task{}
		err = rows.Scan(&singleTask.Id, &singleTask.Nom, &singleTask.Description, &singleTask.Etat, &singleTask.Priorite, &singleTask.Deadline)

		if err != nil {
			return nil, err
		}

		listeTask = append(listeTask, singleTask)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return listeTask, err
}
func GetTasksByDate(userID int, date string) ([]Task, error) {
	parsedDate, err := time.Parse("02012006", date)
	if err != nil {
		return nil, err
	}
	formattedDate := parsedDate.Format("02/01/2006")

	query := "SELECT T.idTache, T.idTodoList, T.nom, T.description, T.etat, T.priorite, T.deadline " +
		"FROM TACHE T INNER JOIN TODOLISTS L ON T.idTodoList = L.idTodoLists " +
		"WHERE T.deadline = ? AND L.idUser = ?"

	rows, err := DB.Query(query, formattedDate, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	taskList := make([]Task, 0)
	for rows.Next() {
		singleTask := Task{}
		err := rows.Scan(&singleTask.Id, &singleTask.IdTodoList, &singleTask.Nom, &singleTask.Description, &singleTask.Etat, &singleTask.Priorite, &singleTask.Deadline)
		if err != nil {
			return nil, err
		}
		taskList = append(taskList, singleTask)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return taskList, nil
}

func AddTodoList(todo Todolist, idUser int) error {
	_, err := DB.Exec("INSERT INTO TODOLISTS (nom, idUser) VALUES (?, ?)", todo.Nom, idUser)
	return err
}
func AddTask(task Task) error {
	fmt.Printf("Adding task: %+v\n", task)
	_, err := DB.Exec("INSERT INTO TACHE (idTodoList, nom, description, etat, priorite, deadline) VALUES (?, ?, ?, ?, ?, ?)",
		task.IdTodoList, task.Nom, task.Description, task.Etat, task.Priorite, task.Deadline)
	return err
}
func DeleteTasksForTodoList(todoListID int) error {
	query := "DELETE FROM TACHE WHERE idTodoList = ?"

	result, err := DB.Exec(query, todoListID)
	if err != nil {
		return fmt.Errorf("error deleting tasks for todolist: %v", err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error retrieving the number of affected rows: %v", err)
	}

	return nil
}

func DeleteTodoList(id int) error {
	if err := DeleteTasksForTodoList(id); err != nil {
		return err
	}
	query := "DELETE FROM TODOLISTS WHERE idTodoLists = ?"

	result, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression de la liste de tâches: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération du nombre de lignes affectées: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("aucune liste de tâches trouvée avec l'ID %d", id)
	}

	return nil
}

func UpdateTask(id int, updatedTask Task) error {
	query := "UPDATE TACHE SET nom=?, description=?, etat=?, priorite=?, deadline=? WHERE idTache=?"

	_, err := DB.Exec(query, updatedTask.Nom, updatedTask.Description, updatedTask.Etat, updatedTask.Priorite, updatedTask.Deadline, id)
	if err != nil {
		return fmt.Errorf("erreur lors de la mise à jour de la tâche: %v", err)
	}

	return nil
}

func UpdateTodoList(id int, updatedTodoList Todolist) error {
	query := "UPDATE TODOLISTS SET nom=? WHERE idTodoLists=?"

	_, err := DB.Exec(query, updatedTodoList.Nom, id)
	if err != nil {
		return fmt.Errorf("erreur lors de la mise à jour de la liste de tâches: %v", err)
	}

	return nil
}

func DeleteTask(idTodoList, idTask int) error {
	query := "DELETE FROM TACHE WHERE idTodoList = ? AND idTache = ?"

	result, err := DB.Exec(query, idTodoList, idTask)
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression de la tâche: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération du nombre de lignes affectées: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("aucune tâche trouvée avec l'ID %d dans la liste de tâches avec l'ID %d", idTask, idTodoList)
	}

	return nil
}

func DeleteTaskById(idTask int) error {
	query := "DELETE FROM TACHE WHERE idTache = ?"

	result, err := DB.Exec(query, idTask)
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression de la tâche: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération du nombre de lignes affectées : %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("aucune tâche trouvée avec l'ID %d", idTask)
	}

	return nil
}

func GetAllUrgentTasks(userID int) ([]Task, error) {
	priority := 3
	etat := 0
	query := "SELECT T.idTache, T.idTodoList, T.nom, T.description, T.etat, T.priorite, T.deadline " +
		"FROM TACHE T INNER JOIN TODOLISTS L ON T.idTodoList = L.idTodoLists " +
		"WHERE T.priorite = ? AND T.etat = ? AND L.idUser = ?"

	rows, err := DB.Query(query, priority, etat, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	listeTask := make([]Task, 0)

	for rows.Next() {
		singleTask := Task{}
		err = rows.Scan(&singleTask.Id, &singleTask.IdTodoList, &singleTask.Nom, &singleTask.Description, &singleTask.Etat, &singleTask.Priorite, &singleTask.Deadline)

		if err != nil {
			return nil, err
		}

		listeTask = append(listeTask, singleTask)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return listeTask, nil
}

func UpdateTaskById(id int, updatedTask Task) error {
	query := "UPDATE TACHE SET nom=?, description=?, etat=?, priorite=?, deadline=? WHERE idTache=?"

	_, err := DB.Exec(query, updatedTask.Nom, updatedTask.Description, updatedTask.Etat, updatedTask.Priorite, updatedTask.Deadline, id)
	if err != nil {
		return fmt.Errorf("error updating task by ID: %v", err)
	}

	return nil
}
