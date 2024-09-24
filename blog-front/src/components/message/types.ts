export type NoticeType = 'success' | 'waring' | 'info' | 'danger' | 'loading';

export interface Notice {
  text: string;
  key: string;
  type: NoticeType;
}

export interface NoticeResult {
  key: string;
  close: () => void;
}
