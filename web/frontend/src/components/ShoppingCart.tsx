import React, { useEffect, useState } from 'react';
import { Product, Cart } from '../types';

const DEMO_PRODUCTS = [
  {
    id: "1",
    name: "Gaming Laptop",
    price: 1299.99,
    stock: 10,
    emoji: "💻"
  },
  {
    id: "2",
    name: "Wireless Headphones",
    price: 199.99,
    stock: 15,
    emoji: "🎧"
  },
  {
    id: "3",
    name: "Smartphone",
    price: 899.99,
    stock: 8,
    emoji: "📱"
  },
  {
    id: "4",
    name: "Smart Watch",
    price: 299.99,
    stock: 12,
    emoji: "⌚"
  },
  {
    id: "5",
    name: "Tablet",
    price: 499.99,
    stock: 6,
    emoji: "📱"
  },
  {
    id: "6",
    name: "Wireless Mouse",
    price: 49.99,
    stock: 20,
    emoji: "🖱️"
  }
];

const ShoppingCart = () => {
  const [products, setProducts] = useState<Product[]>([]);
  const [cart, setCart] = useState<Cart | null>(null);
  const [loading, setLoading] = useState(true);
  const userId = "user123";

  const fetchProducts = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/v1/products');
      const data = await response.json();
      // 如果API沒有返回數據，使用演示數據
      setProducts(data.length > 0 ? data : DEMO_PRODUCTS);
    } catch (error) {
      console.error('Error fetching products:', error);
      // 如果API調用失敗，使用演示數據
      setProducts(DEMO_PRODUCTS);
    }
  };

  const fetchCart = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/v1/cart/${userId}`);
      const data = await response.json();
      setCart(data);
    } catch (error) {
      console.error('Error fetching cart:', error);
    }
  };

  const addToCart = async (productId: string) => {
    try {
      await fetch('http://localhost:8080/api/v1/cart', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          user_id: userId,
          product_id: productId,
          quantity: 1,
        }),
      });
      fetchCart();
    } catch (error) {
      console.error('Error adding to cart:', error);
    }
  };

  const updateQuantity = async (productId: string, quantity: number) => {
    try {
      await fetch(`http://localhost:8080/api/v1/cart/${userId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          product_id: productId,
          quantity,
        }),
      });
      fetchCart();
    } catch (error) {
      console.error('Error updating quantity:', error);
    }
  };

  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      await Promise.all([fetchProducts(), fetchCart()]);
      setLoading(false);
    };
    loadData();
  }, []);

  if (loading) {
    return (
      <div className="fixed inset-0 flex items-center justify-center bg-white bg-opacity-75 backdrop-blur-sm">
        <div className="animate-spin rounded-full h-16 w-16 border-t-4 border-b-4 border-indigo-600"></div>
      </div>
    );
  }

  return (
    
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow-lg sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16 items-center">
            <span className="text-2xl font-bold text-indigo-600">
              Tech Store 🛍️
            </span>
            <div className="relative">
              <span className="text-2xl">🛒</span>
              {cart && cart.items.length > 0 && (
                <span className="absolute -top-2 -right-2 bg-red-500 text-white text-xs w-5 h-5 flex items-center justify-center rounded-full">
                  {cart.items.length}
                </span>
              )}
            </div>
          </div>
        </div>
      </nav>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <div className="lg:col-span-2">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {products.map((product) => (
                <div key={product.id} className="bg-white rounded-lg shadow-sm hover:shadow-md transition-all duration-300 transform hover:-translate-y-1">
                  <div className="p-6">
                    <div className="text-center mb-4">
                      <span className="text-6xl" role="img" aria-label={product.name}>
                        {product.emoji || "📦"}
                      </span>
                    </div>
                    <h3 className="text-lg font-semibold text-gray-800 text-center mb-2">{product.name}</h3>
                    <div className="mt-4">
                      <div className="flex justify-between items-center mb-4">
                        <span className="text-lg font-bold text-indigo-600">${product.price.toFixed(2)}</span>
                        <span className="text-sm text-gray-500">Stock: {product.stock}</span>
                      </div>
                      <button
                        onClick={() => addToCart(product.id)}
                        className="w-full bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700 
                          transition-colors duration-200 flex items-center justify-center gap-2"
                      >
                        <span>Add to Cart</span>
                        <span>+</span>
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Cart Section */}
          <div className="lg:sticky lg:top-24">
            <div className="bg-white rounded-lg shadow-sm p-6">
              <h2 className="text-xl font-semibold mb-6 flex items-center gap-2">
                Shopping Cart
              </h2>

              {cart && cart.items.length > 0 ? (
                <div className="space-y-4">
                  {cart.items.map((item) => (
                    <div key={item.product_id} className="flex flex-col gap-2 py-4 border-b border-gray-100">
                      <div className="flex justify-between items-start">
                        <div>
                          <h3 className="font-medium text-gray-800">{item.name}</h3>
                          <p className="text-sm text-gray-500">${item.price.toFixed(2)} each</p>
                        </div>
                        <div className="flex items-center space-x-2">
                          <button
                            onClick={() => updateQuantity(item.product_id, Math.max(0, item.quantity - 1))}
                            className="w-8 h-8 flex items-center justify-center rounded-full bg-gray-100 
                              hover:bg-gray-200 transition-colors text-gray-600"
                          >
                            -
                          </button>
                          <span className="w-8 text-center">{item.quantity}</span>
                          <button
                            onClick={() => updateQuantity(item.product_id, item.quantity + 1)}
                            className="w-8 h-8 flex items-center justify-center rounded-full bg-gray-100 
                              hover:bg-gray-200 transition-colors text-gray-600"
                          >
                            +
                          </button>
                          <button
                            onClick={() => updateQuantity(item.product_id, 0)}
                            className="w-8 h-8 flex items-center justify-center rounded-full bg-red-100 
                              hover:bg-red-200 transition-colors text-red-600"
                          >
                            ×
                          </button>
                        </div>
                      </div>
                      <div className="text-right text-sm text-gray-600">
                        Subtotal: ${(item.price * item.quantity).toFixed(2)}
                      </div>
                    </div>
                  ))}

                  <div className="pt-4">
                    <div className="flex justify-between items-center text-lg font-semibold">
                      <span>Total</span>
                      <span className="text-indigo-600">${cart.total.toFixed(2)}</span>
                    </div>
                    <button className="mt-6 w-full bg-gradient-to-r from-indigo-600 to-purple-600 text-white py-3 
                      rounded-lg hover:from-indigo-700 hover:to-purple-700 transition-all duration-200 
                      focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 font-medium">
                      Proceed to Checkout
                    </button>
                  </div>
                </div>
              ) : (
                <div className="text-center py-12">
                  <span className="text-6xl mb-4 block">🛒</span>
                  <p className="text-gray-500">Your cart is empty</p>
                  <p className="text-sm text-gray-400 mt-2">Add some products to your cart</p>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ShoppingCart;