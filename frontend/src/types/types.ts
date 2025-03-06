export type Product = {
    id: number,
    name: string,
    brand: string, 
    description: string,
    price: number,
    count: number
}

export type Category = {
    id: number,
    name: string,
    description: string
}

export type ErrorObject = {
    error: number | null,
    message: string | null
}