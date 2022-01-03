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
    image_path?: string,
    comments?: IPost[]

}

export interface ICategory {
    id: number;
    name: string;
}
