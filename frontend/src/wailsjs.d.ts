interface Window {
  go: {
    main: {
      App: {
        SearchBooks(keyword: string): Promise<Book[]>;
        DownloadBook(bookId: string, format: string): Promise<void>;
        GetConfig(): Promise<Config>;
        SaveConfig(config: Config): Promise<void>;
  
      };
    };
  };
} 