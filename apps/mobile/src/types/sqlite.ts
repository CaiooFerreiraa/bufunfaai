export interface LocalSnapshot<TData> {
  readonly key: string;
  readonly payload: TData;
  readonly syncedAt: string;
}
