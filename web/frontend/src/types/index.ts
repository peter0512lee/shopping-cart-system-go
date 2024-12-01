// src/types/index.ts
export interface Product {
  id: string;
  name: string;
  price: number;
  stock: number;
  emoji: string;
}

export interface CartItem {
  product_id: string;
  name: string;
  price: number;
  quantity: number;
}

export interface Cart {
  user_id: string;
  items: CartItem[];
  total: number;
}