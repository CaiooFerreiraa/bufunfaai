import * as SQLite from 'expo-sqlite';

import type { LocalSnapshot } from '@/types/sqlite';

const DATABASE_NAME: string = 'bufunfaai.db';

let databasePromise: Promise<SQLite.SQLiteDatabase> | null = null;

async function getDatabase(): Promise<SQLite.SQLiteDatabase> {
  if (!databasePromise) {
    databasePromise = SQLite.openDatabaseAsync(DATABASE_NAME);
  }

  return databasePromise;
}

export async function initializeSQLite(): Promise<void> {
  const database: SQLite.SQLiteDatabase = await getDatabase();
  await database.execAsync(`
    CREATE TABLE IF NOT EXISTS app_snapshots (
      key TEXT PRIMARY KEY NOT NULL,
      payload TEXT NOT NULL,
      synced_at TEXT NOT NULL
    );
  `);
}

export async function saveSnapshot<TData>(snapshot: LocalSnapshot<TData>): Promise<void> {
  const database: SQLite.SQLiteDatabase = await getDatabase();
  await database.runAsync(
    `INSERT OR REPLACE INTO app_snapshots (key, payload, synced_at) VALUES (?, ?, ?)`,
    snapshot.key,
    JSON.stringify(snapshot.payload),
    snapshot.syncedAt,
  );
}

export async function getSnapshot<TData>(key: string): Promise<LocalSnapshot<TData> | null> {
  const database: SQLite.SQLiteDatabase = await getDatabase();
  const row: { readonly key: string; readonly payload: string; readonly synced_at: string } | null =
    await database.getFirstAsync(
      `SELECT key, payload, synced_at FROM app_snapshots WHERE key = ?`,
      key,
    );

  if (!row) {
    return null;
  }

  return {
    key: row.key,
    payload: JSON.parse(row.payload) as TData,
    syncedAt: row.synced_at,
  };
}
