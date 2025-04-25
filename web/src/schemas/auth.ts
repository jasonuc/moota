import { z } from "zod";

export const registerFormSchema = z.object({
    username: z.string()
        .min(3, "Name must be at least 2 characters")
        .max(50, "Name cannot exceed 50 characters")
        .trim(),
    email: z.string()
        .email("Please enter a valid email address")
        .trim()
        .transform(val => val.toLowerCase()),
    password: z.string()
        .min(8, "Password must be at least 8 characters")
        .max(72, "Password cannot exceed 72 characters")
        .regex(/[A-Z]/, "Password must contain at least one uppercase letter")
        .regex(/[0-9]/, "Password must contain at least one number"),
    confirmPassword: z.string()
})