import { z } from "zod";

export const changeUsernameFormSchema = z.object({
  newUsername: z
    .string()
    .min(3, "Name must be at least 3 characters")
    .max(30, "Name cannot exceed 30 characters")
    .trim(),
});

export const changeEmailFormSchema = z.object({
  newEmail: z.string().email("Please enter a valid email address").trim(),
});

export const changePasswordFormSchema = z.object({
  oldPassword: z.string().min(1, "Password is required"),
  newPassword: z
    .string()
    .min(8, "Password must be at least 8 characters")
    .max(72, "Password cannot exceed 72 characters")
    .regex(/[A-Z]/, "Password must contain at least one uppercase letter")
    .regex(/[0-9]/, "Password must contain at least one number"),
  confirmNewPassword: z.string().min(1, "Password is required"),
});
