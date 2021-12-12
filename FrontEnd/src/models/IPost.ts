export interface IPost {
    id: number;
    user_id: string;
    content: string;
    created_at: string;
    subject: string;
    user_login: string;
    likes: number;
    dislikes: number;
    categories: string;
    "parent_id": number;
}

export interface ICategory {
    id: number;
    name: string;
}
