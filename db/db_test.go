package db

import (
	"context"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"

	"github.com/0xPolygon/cdk-data-availability/types"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	constructorExpect(mock)

	wdb := sqlx.NewDb(db, "postgres")

	_, err = New(context.Background(), wdb)
	require.NoError(t, err)
}

func Test_DB_StoreLastProcessedBlock(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name      string
		task      string
		block     uint64
		returnErr error
	}{
		{
			name:  "value inserted",
			task:  "task1",
			block: 1,
		},
		{
			name:      "error returned",
			task:      "task1",
			block:     1,
			returnErr: errors.New("test error"),
		},
	}

	for _, tt := range testTable {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)

			defer db.Close()

			constructorExpect(mock)

			expected := mock.ExpectExec(`INSERT INTO data_node\.sync_tasks \(task, block\) VALUES \(\$1, \$2\) ON CONFLICT \(task\) DO UPDATE SET block = EXCLUDED\.block, processed = NOW\(\)`).
				WithArgs(tt.task, tt.block)
			if tt.returnErr != nil {
				expected.WillReturnError(tt.returnErr)
			} else {
				expected.WillReturnResult(sqlmock.NewResult(1, 1))
			}

			wdb := sqlx.NewDb(db, "postgres")

			dbPG, err := New(context.Background(), wdb)
			require.NoError(t, err)

			err = dbPG.StoreLastProcessedBlock(context.Background(), tt.block, tt.task)
			if tt.returnErr != nil {
				require.ErrorIs(t, err, tt.returnErr)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_DB_GetLastProcessedBlock(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name      string
		task      string
		block     uint64
		returnErr error
	}{
		{
			name:  "successfully selected block",
			task:  "task1",
			block: 1,
		},
		{
			name:      "error returned",
			task:      "task1",
			block:     1,
			returnErr: errors.New("test error"),
		},
	}

	for _, tt := range testTable {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)

			defer db.Close()

			constructorExpect(mock)

			mock.ExpectExec(`INSERT INTO data_node\.sync_tasks \(task, block\) VALUES \(\$1, \$2\) ON CONFLICT \(task\) DO UPDATE SET block = EXCLUDED\.block, processed = NOW\(\)`).
				WithArgs(tt.task, tt.block).
				WillReturnResult(sqlmock.NewResult(1, 1))

			expected := mock.ExpectQuery(`SELECT block FROM data_node\.sync_tasks WHERE task = \$1`).
				WithArgs(tt.task)

			if tt.returnErr != nil {
				expected.WillReturnError(tt.returnErr)
			} else {
				expected.WillReturnRows(sqlmock.NewRows([]string{"block"}).AddRow(tt.block))
			}

			wdb := sqlx.NewDb(db, "postgres")

			dbPG, err := New(context.Background(), wdb)
			require.NoError(t, err)

			err = dbPG.StoreLastProcessedBlock(context.Background(), tt.block, tt.task)
			require.NoError(t, err)

			actual, err := dbPG.GetLastProcessedBlock(context.Background(), tt.task)
			if tt.returnErr != nil {
				require.ErrorIs(t, err, tt.returnErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.block, actual)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_DB_StoreUnresolvedBatchKeys(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name      string
		bk        []types.BatchKey
		returnErr error
	}{
		{
			name: "no values inserted",
		},
		{
			name: "one value inserted",
			bk: []types.BatchKey{{
				Number: 1,
				Hash:   common.HexToHash("key1"),
			}},
		},
		{
			name: "several values inserted",
			bk: []types.BatchKey{{
				Number: 1,
				Hash:   common.HexToHash("key1"),
			}, {
				Number: 2,
				Hash:   common.HexToHash("key2"),
			}},
		},
		{
			name: "error returned",
			bk: []types.BatchKey{{
				Number: 1,
				Hash:   common.HexToHash("key1"),
			}},
			returnErr: errors.New("test error"),
		},
	}

	for _, tt := range testTable {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)

			wdb := sqlx.NewDb(db, "postgres")

			mock.ExpectPrepare(regexp.QuoteMeta(storeLastProcessedBlockSQL))
			mock.ExpectPrepare(regexp.QuoteMeta(getLastProcessedBlockSQL))
			mock.ExpectPrepare(regexp.QuoteMeta(getUnresolvedBatchKeysSQL))
			mock.ExpectPrepare(regexp.QuoteMeta(getOffchainDataSQL))
			mock.ExpectPrepare(regexp.QuoteMeta(countOffchainDataSQL))
			mock.ExpectPrepare(regexp.QuoteMeta(selectOffchainDataGapsSQL))

			dbPG, err := New(context.Background(), wdb)
			require.NoError(t, err)

			defer db.Close()

			mock.ExpectBegin()
			mock.ExpectPrepare(regexp.QuoteMeta(storeUnresolvedBatchesSQL))
			for _, o := range tt.bk {
				expected := mock.ExpectExec(regexp.QuoteMeta(storeUnresolvedBatchesSQL)).
					WithArgs(o.Number, o.Hash.Hex())
				if tt.returnErr != nil {
					expected.WillReturnError(tt.returnErr)
				} else {
					expected.WillReturnResult(sqlmock.NewResult(int64(len(tt.bk)), int64(len(tt.bk))))
				}
			}
			if tt.returnErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			err = dbPG.StoreUnresolvedBatchKeys(context.Background(), tt.bk)
			if tt.returnErr != nil {
				require.ErrorIs(t, err, tt.returnErr)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_DB_GetUnresolvedBatchKeys(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name      string
		bks       []types.BatchKey
		returnErr error
	}{
		{
			name: "successfully selected data",
			bks: []types.BatchKey{{
				Number: 1,
				Hash:   common.HexToHash("key1"),
			}},
		},
		{
			name: "error returned",
			bks: []types.BatchKey{{
				Number: 1,
				Hash:   common.HexToHash("key1"),
			}},
			returnErr: errors.New("test error"),
		},
	}

	for _, tt := range testTable {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)

			defer db.Close()

			constructorExpect(mock)

			wdb := sqlx.NewDb(db, "postgres")
			dbPG, err := New(context.Background(), wdb)
			require.NoError(t, err)

			// Seed data
			seedUnresolvedBatchKeys(t, dbPG, mock, tt.bks)

			var limit = uint(10)
			expected := mock.ExpectQuery(`SELECT num, hash FROM data_node\.unresolved_batches LIMIT \$1\;`).WithArgs(limit)

			if tt.returnErr != nil {
				expected.WillReturnError(tt.returnErr)
			} else {
				for _, bk := range tt.bks {
					expected.WillReturnRows(sqlmock.NewRows([]string{"num", "hash"}).AddRow(bk.Number, bk.Hash.Hex()))
				}
			}

			data, err := dbPG.GetUnresolvedBatchKeys(context.Background(), limit)
			if tt.returnErr != nil {
				require.ErrorIs(t, err, tt.returnErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.bks, data)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_DB_DeleteUnresolvedBatchKeys(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name      string
		bks       []types.BatchKey
		returnErr error
	}{
		{
			name: "value deleted",
			bks: []types.BatchKey{{
				Number: 1,
				Hash:   common.HexToHash("key1"),
			}},
		},
		{
			name: "error returned",
			bks: []types.BatchKey{{
				Number: 1,
				Hash:   common.HexToHash("key1"),
			}},
			returnErr: errors.New("test error"),
		},
	}

	for _, tt := range testTable {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)

			defer db.Close()

			constructorExpect(mock)

			mock.ExpectBegin()
			mock.ExpectPrepare(regexp.QuoteMeta(deleteUnresolvedBatchKeysSQL))
			for _, bk := range tt.bks {
				expected := mock.ExpectExec(regexp.QuoteMeta(deleteUnresolvedBatchKeysSQL)).
					WithArgs(bk.Number, bk.Hash.Hex())
				if tt.returnErr != nil {
					expected.WillReturnError(tt.returnErr)
				} else {
					expected.WillReturnResult(sqlmock.NewResult(int64(len(tt.bks)), int64(len(tt.bks))))
				}
			}
			if tt.returnErr != nil {
				mock.ExpectRollback()
			} else {
				mock.ExpectCommit()
			}

			wdb := sqlx.NewDb(db, "postgres")

			dbPG, err := New(context.Background(), wdb)
			require.NoError(t, err)

			err = dbPG.DeleteUnresolvedBatchKeys(context.Background(), tt.bks)
			if tt.returnErr != nil {
				require.ErrorIs(t, err, tt.returnErr)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_DB_StoreOffChainData(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name      string
		od        []types.OffChainData
		returnErr error
	}{
		{
			name: "no values inserted",
		},
		{
			name: "one value inserted",
			od: []types.OffChainData{{
				Key:   common.HexToHash("key1"),
				Value: []byte("value1"),
			}},
		},
		{
			name: "several values inserted",
			od: []types.OffChainData{{
				Key:   common.HexToHash("key1"),
				Value: []byte("value1"),
			}, {
				Key:   common.HexToHash("key2"),
				Value: []byte("value2"),
			}},
		},
		{
			name: "error returned",
			od: []types.OffChainData{{
				Key:   common.HexToHash("key1"),
				Value: []byte("value1"),
			}},
			returnErr: errors.New("test error"),
		},
	}

	for _, tt := range testTable {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)

			constructorExpect(mock)

			wdb := sqlx.NewDb(db, "postgres")
			dbPG, err := New(context.Background(), wdb)
			require.NoError(t, err)

			defer db.Close()

			mock.ExpectBegin()
			mock.ExpectPrepare(regexp.QuoteMeta(storeOffChainDataSQL))
			for _, o := range tt.od {
				expected := mock.ExpectExec(regexp.QuoteMeta(storeOffChainDataSQL)).
					WithArgs(o.Key.Hex(), common.Bytes2Hex(o.Value), o.BatchNum)
				if tt.returnErr != nil {
					expected.WillReturnError(tt.returnErr)
				} else {
					expected.WillReturnResult(sqlmock.NewResult(int64(len(tt.od)), int64(len(tt.od))))
				}
			}
			if tt.returnErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			err = dbPG.StoreOffChainData(context.Background(), tt.od)
			if tt.returnErr != nil {
				require.ErrorIs(t, err, tt.returnErr)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_DB_GetOffChainData(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name      string
		od        []types.OffChainData
		key       common.Hash
		expected  *types.OffChainData
		returnErr error
	}{
		{
			name: "successfully selected value",
			od: []types.OffChainData{{
				Key:      common.HexToHash("key1"),
				Value:    []byte("value1"),
				BatchNum: 1,
			}},
			key: common.BytesToHash([]byte("key1")),
			expected: &types.OffChainData{
				Key:      common.HexToHash("key1"),
				Value:    []byte("value1"),
				BatchNum: 1,
			},
		},
		{
			name: "error returned",
			od: []types.OffChainData{{
				Key:   common.HexToHash("key1"),
				Value: []byte("value1"),
			}},
			key:       common.BytesToHash([]byte("key1")),
			returnErr: errors.New("test error"),
		},
		{
			name: "no rows",
			od: []types.OffChainData{{
				Key:   common.HexToHash("key1"),
				Value: []byte("value1"),
			}},
			key:       common.BytesToHash([]byte("undefined")),
			returnErr: ErrStateNotSynchronized,
		},
	}

	for _, tt := range testTable {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)

			constructorExpect(mock)

			wdb := sqlx.NewDb(db, "postgres")
			dbPG, err := New(context.Background(), wdb)
			require.NoError(t, err)

			defer db.Close()

			// Seed data
			seedOffchainData(t, dbPG, mock, tt.od)

			expected := mock.ExpectQuery(regexp.QuoteMeta(getOffchainDataSQL)).
				WithArgs(tt.key.Hex())

			if tt.returnErr != nil {
				expected.WillReturnError(tt.returnErr)
			} else {
				expected.WillReturnRows(sqlmock.NewRows([]string{"key", "value", "batch_num"}).
					AddRow(tt.expected.Key.Hex(), common.Bytes2Hex(tt.expected.Value), tt.expected.BatchNum))
			}

			data, err := dbPG.GetOffChainData(context.Background(), tt.key)
			if tt.returnErr != nil {
				require.ErrorIs(t, err, tt.returnErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, data)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_DB_ListOffChainData(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name      string
		od        []types.OffChainData
		keys      []common.Hash
		expected  []types.OffChainData
		sql       string
		returnErr error
	}{
		{
			name: "successfully selected one value",
			od: []types.OffChainData{{
				Key:   common.HexToHash("key1"),
				Value: []byte("value1"),
			}},
			keys: []common.Hash{
				common.BytesToHash([]byte("key1")),
			},
			expected: []types.OffChainData{
				{
					Key:      common.BytesToHash([]byte("key1")),
					Value:    []byte("value1"),
					BatchNum: 0,
				},
			},
			sql: `SELECT key, value, batch_num FROM data_node\.offchain_data WHERE key IN \(\$1\)`,
		},
		{
			name: "successfully selected two values",
			od: []types.OffChainData{{
				Key:      common.HexToHash("key1"),
				Value:    []byte("value1"),
				BatchNum: 1,
			}, {
				Key:      common.HexToHash("key2"),
				Value:    []byte("value2"),
				BatchNum: 2,
			}},
			keys: []common.Hash{
				common.BytesToHash([]byte("key1")),
				common.BytesToHash([]byte("key2")),
			},
			expected: []types.OffChainData{
				{
					Key:      common.BytesToHash([]byte("key1")),
					Value:    []byte("value1"),
					BatchNum: 1,
				},
				{
					Key:      common.BytesToHash([]byte("key2")),
					Value:    []byte("value2"),
					BatchNum: 2,
				},
			},
			sql: `SELECT key, value, batch_num FROM data_node\.offchain_data WHERE key IN \(\$1\, \$2\)`,
		},
		{
			name: "error returned",
			od: []types.OffChainData{{
				Key:   common.HexToHash("key1"),
				Value: []byte("value1"),
			}},
			keys: []common.Hash{
				common.BytesToHash([]byte("key1")),
			},
			sql:       `SELECT key, value, batch_num FROM data_node\.offchain_data WHERE key IN \(\$1\)`,
			returnErr: errors.New("test error"),
		},
		{
			name: "no rows",
			od: []types.OffChainData{{
				Key:   common.HexToHash("key1"),
				Value: []byte("value1"),
			}},
			keys: []common.Hash{
				common.BytesToHash([]byte("undefined")),
			},
			sql:       `SELECT key, value, batch_num FROM data_node\.offchain_data WHERE key IN \(\$1\)`,
			returnErr: ErrStateNotSynchronized,
		},
	}

	for _, tt := range testTable {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)

			defer db.Close()

			constructorExpect(mock)

			wdb := sqlx.NewDb(db, "postgres")
			dbPG, err := New(context.Background(), wdb)
			require.NoError(t, err)

			// Seed data
			seedOffchainData(t, dbPG, mock, tt.od)

			preparedKeys := make([]driver.Value, len(tt.keys))
			for i, key := range tt.keys {
				preparedKeys[i] = key.Hex()
			}

			expected := mock.ExpectQuery(tt.sql).
				WithArgs(preparedKeys...)

			if tt.returnErr != nil {
				expected.WillReturnError(tt.returnErr)
			} else {
				returnData := sqlmock.NewRows([]string{"key", "value", "batch_num"})

				for _, data := range tt.expected {
					returnData = returnData.AddRow(data.Key.Hex(), common.Bytes2Hex(data.Value), data.BatchNum)
				}

				expected.WillReturnRows(returnData)
			}

			data, err := dbPG.ListOffChainData(context.Background(), tt.keys)
			if tt.returnErr != nil {
				require.ErrorIs(t, err, tt.returnErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, data)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_DB_CountOffchainData(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name      string
		od        []types.OffChainData
		count     uint64
		returnErr error
	}{
		{
			name: "two values found",
			od: []types.OffChainData{{
				Key:   common.HexToHash("key1"),
				Value: []byte("value1"),
			}, {
				Key:   common.HexToHash("key1"),
				Value: []byte("value2"),
			}},
			count: 2,
		},
		{
			name:  "no values found",
			count: 0,
		},
		{
			name: "error returned",
			od: []types.OffChainData{{
				Key:   common.HexToHash("key1"),
				Value: []byte("value1"),
			}},
			returnErr: errors.New("test error"),
		},
	}

	for _, tt := range testTable {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)

			defer db.Close()

			constructorExpect(mock)

			wdb := sqlx.NewDb(db, "postgres")
			dbPG, err := New(context.Background(), wdb)
			require.NoError(t, err)

			// Seed data
			seedOffchainData(t, dbPG, mock, tt.od)

			expected := mock.ExpectQuery(regexp.QuoteMeta(countOffchainDataSQL))

			if tt.returnErr != nil {
				expected.WillReturnError(tt.returnErr)
			} else {
				expected.WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(tt.count))
			}

			actual, err := dbPG.CountOffchainData(context.Background())
			if tt.returnErr != nil {
				require.ErrorIs(t, err, tt.returnErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.count, actual)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_DB_DetectOffchainDataGaps(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name      string
		seed      []types.OffChainData
		gaps      map[uint64]uint64
		returnErr error
	}{
		{
			name: "one gap found",
			seed: []types.OffChainData{{
				Key:      common.HexToHash("key1"),
				Value:    []byte("value1"),
				BatchNum: 1,
			}, {
				Key:      common.HexToHash("key2"),
				Value:    []byte("value2"),
				BatchNum: 2,
			}, {
				Key:      common.HexToHash("key4"),
				Value:    []byte("value4"),
				BatchNum: 4,
			}},
			gaps: map[uint64]uint64{
				2: 4,
			},
		},
		{
			name: "error returned",
			seed: []types.OffChainData{{
				Key:   common.HexToHash("key1"),
				Value: []byte("value1"),
			}},
			returnErr: errors.New("test error"),
		},
	}

	for _, tt := range testTable {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)

			defer db.Close()

			constructorExpect(mock)

			wdb := sqlx.NewDb(db, "postgres")
			dbPG, err := New(context.Background(), wdb)
			require.NoError(t, err)

			// Seed data
			seedOffchainData(t, dbPG, mock, tt.seed)

			expected := mock.ExpectQuery(regexp.QuoteMeta(selectOffchainDataGapsSQL))

			if tt.returnErr != nil {
				expected.WillReturnError(tt.returnErr)
			} else {
				rows := sqlmock.NewRows([]string{"current_batch_num", "next_batch_num"})
				for k, v := range tt.gaps {
					rows.AddRow(k, v)
				}
				expected.WillReturnRows(rows)
			}

			actual, err := dbPG.DetectOffchainDataGaps(context.Background())
			if tt.returnErr != nil {
				require.ErrorIs(t, err, tt.returnErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.gaps, actual)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func constructorExpect(mock sqlmock.Sqlmock) {
	mock.ExpectPrepare(regexp.QuoteMeta(storeLastProcessedBlockSQL))
	mock.ExpectPrepare(regexp.QuoteMeta(getLastProcessedBlockSQL))
	mock.ExpectPrepare(regexp.QuoteMeta(getUnresolvedBatchKeysSQL))
	mock.ExpectPrepare(regexp.QuoteMeta(getOffchainDataSQL))
	mock.ExpectPrepare(regexp.QuoteMeta(countOffchainDataSQL))
	mock.ExpectPrepare(regexp.QuoteMeta(selectOffchainDataGapsSQL))
}

func seedOffchainData(t *testing.T, db DB, mock sqlmock.Sqlmock, od []types.OffChainData) {
	t.Helper()

	mock.ExpectBegin()
	mock.ExpectPrepare(regexp.QuoteMeta(storeOffChainDataSQL))
	for i, o := range od {
		mock.ExpectExec(regexp.QuoteMeta(storeOffChainDataSQL)).
			WithArgs(o.Key.Hex(), common.Bytes2Hex(o.Value), o.BatchNum).
			WillReturnResult(sqlmock.NewResult(int64(i+1), int64(i+1)))
	}
	mock.ExpectCommit()

	err := db.StoreOffChainData(context.Background(), od)
	require.NoError(t, err)
}

func seedUnresolvedBatchKeys(t *testing.T, db DB, mock sqlmock.Sqlmock, bk []types.BatchKey) {
	t.Helper()

	mock.ExpectBegin()
	mock.ExpectPrepare(regexp.QuoteMeta(storeUnresolvedBatchesSQL))
	for i, o := range bk {
		mock.ExpectExec(regexp.QuoteMeta(storeUnresolvedBatchesSQL)).
			WithArgs(o.Number, o.Hash.Hex()).
			WillReturnResult(sqlmock.NewResult(int64(i+1), int64(i+1)))
	}
	mock.ExpectCommit()

	err := db.StoreUnresolvedBatchKeys(context.Background(), bk)
	require.NoError(t, err)
}
