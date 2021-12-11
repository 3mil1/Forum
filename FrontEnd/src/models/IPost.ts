export interface IPost {
    post: Post;
    likes?: number | null;
    dislikes?: number | null;
}
export interface Post {
    id: number;
    user_id: string;
    content: string;
    created_at: string;
    subject: string;
    parent_id?: number | null;
}
