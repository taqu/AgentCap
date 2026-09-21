package store

import (
	"context"
	"fmt"

	"github.com/taqu/agentcap/internal/gitparse"
)

// SaveGitDiff stores git diff file and hunk metadata for a result.
func (s *Store) SaveGitDiff(resultID string, files []gitparse.GitDiffFile) error {
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, f := range files {
		_, err := tx.ExecContext(context.Background(), `
			INSERT OR REPLACE INTO git_diff_files
			(result_id, file_index, old_path, new_path, status, additions, deletions, binary, old_mode, new_mode, similarity, raw_start, raw_end)
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			resultID, f.Index, f.OldPath, f.NewPath, f.Status,
			f.Additions, f.Deletions, boolToInt(f.Binary),
			f.OldMode, f.NewMode, f.Similarity, f.RawStart, f.RawEnd)
		if err != nil {
			return fmt.Errorf("store: save git file: %w", err)
		}

		for _, h := range f.Hunks {
			_, err := tx.ExecContext(context.Background(), `
				INSERT OR REPLACE INTO git_diff_hunks
				(result_id, file_index, hunk_index, old_start, old_lines, new_start, new_lines, header, raw_start, raw_end)
				VALUES (?,?,?,?,?,?,?,?,?,?)`,
				resultID, f.Index, h.Index, h.OldStart, h.OldLines,
				h.NewStart, h.NewLines, h.Header, h.RawStart, h.RawEnd)
			if err != nil {
				return fmt.Errorf("store: save git hunk: %w", err)
			}
		}
	}
	return tx.Commit()
}

// GetGitDiffFiles retrieves file metadata for a git diff result.
func (s *Store) GetGitDiffFiles(resultID string) ([]gitparse.GitDiffFile, error) {
	rows, err := s.db.QueryContext(context.Background(), `
		SELECT file_index, old_path, new_path, status, additions, deletions, binary, old_mode, new_mode, similarity, raw_start, raw_end
		FROM git_diff_files WHERE result_id = ? ORDER BY file_index`, resultID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []gitparse.GitDiffFile
	for rows.Next() {
		var f gitparse.GitDiffFile
		var binaryInt int
		err := rows.Scan(&f.Index, &f.OldPath, &f.NewPath, &f.Status,
			&f.Additions, &f.Deletions, &binaryInt,
			&f.OldMode, &f.NewMode, &f.Similarity, &f.RawStart, &f.RawEnd)
		if err != nil {
			continue
		}
		f.Binary = binaryInt != 0
		files = append(files, f)
	}
	return files, nil
}

// GetGitDiffHunks retrieves hunk metadata for a specific file in a git diff result.
func (s *Store) GetGitDiffHunks(resultID string, fileIndex int) ([]gitparse.GitDiffHunk, error) {
	rows, err := s.db.QueryContext(context.Background(), `
		SELECT hunk_index, old_start, old_lines, new_start, new_lines, header, raw_start, raw_end
		FROM git_diff_hunks WHERE result_id = ? AND file_index = ? ORDER BY hunk_index`,
		resultID, fileIndex)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hunks []gitparse.GitDiffHunk
	for rows.Next() {
		var h gitparse.GitDiffHunk
		if err := rows.Scan(&h.Index, &h.OldStart, &h.OldLines, &h.NewStart, &h.NewLines, &h.Header, &h.RawStart, &h.RawEnd); err == nil {
			hunks = append(hunks, h)
		}
	}
	return hunks, nil
}
