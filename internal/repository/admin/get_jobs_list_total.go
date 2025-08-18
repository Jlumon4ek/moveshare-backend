package admin

import (
	"context"
	"fmt"
	"strings"
)

func (r *repository) GetJobsListTotal(ctx context.Context, statuses []string) (int, error) {
	var query string
	var args []interface{}
	
	if len(statuses) == 0 {
		query = `
			SELECT COUNT(j.id)
			FROM jobs j;
		`
		args = []interface{}{}
	} else {
		placeholders := make([]string, len(statuses))
		args = make([]interface{}, len(statuses))
		
		for i, status := range statuses {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
			args[i] = status
		}
		
		query = fmt.Sprintf(`
			SELECT COUNT(j.id)
			FROM jobs j
			WHERE j.job_status IN (%s);
		`, strings.Join(placeholders, ","))
	}

	var count int
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}