export interface Paste {
    PasteName: string;
    Expire: number;
    Privacy: string;
    IsEncrypted: number;
    ReadCount: number;
    ReadLast: number;
    BurnAfter: number;
    Content: string;
    UrlRedirect: number;
    Syntax: string;
    Password: string;
    Files: File[];
    CreatedAt: number;
    UpdatedAt: number;
}

export interface File {
    ID: number;
    Name: string;
    Size: number;
    ContentType: string;
    Blob: any;
}

export interface Config {
    Salt: string;
    Port: string;
    DataDir: string;
    AdminPassword: string;
    PublicList: boolean;
    FileUpload: boolean;
    MaxFileSize: number;
    MaxEncryptionSize: number;
    MaxContentLength: number;
    UploadingPassword: string;
    EternalPaste: boolean;
    ReadCount: boolean;
    BurnAfter: boolean;
    DefaultExpiry: string;
    ShortPasteNames: boolean;
    ShortenRedirectPastes: boolean;
    CountFileUsage: boolean;
}
