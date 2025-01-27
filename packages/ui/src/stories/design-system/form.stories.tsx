import { zodResolver } from "@hookform/resolvers/zod";
import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import { useForm } from "react-hook-form";
import * as z from "zod";
import { Button } from "../../components/button.js";
import { Checkbox } from "../../components/checkbox.js";
import {
	Form,
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "../../components/form.js";
import { Input } from "../../components/input.js";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "../../components/select.js";

const meta = {
	title: "Design System/Form",
	component: Form,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
} satisfies Meta<typeof Form>;

export default meta;
type Story = StoryObj<typeof Form>;

const formSchema = z.object({
	username: z.string().min(2, {
		message: "Username must be at least 2 characters.",
	}),
	email: z.string().email({
		message: "Please enter a valid email address.",
	}),
	type: z.enum(["personal", "business"], {
		required_error: "Please select an account type.",
	}),
	marketing: z.boolean().default(false).optional(),
});

const SignUpForm = () => {
	const form = useForm<z.infer<typeof formSchema>>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			username: "",
			email: "",
			marketing: false,
		},
	});

	function onSubmit(values: z.infer<typeof formSchema>) {}

	return (
		<Form {...form}>
			<form
				onSubmit={form.handleSubmit(onSubmit)}
				className="space-y-6 w-[400px]"
			>
				<FormField
					control={form.control}
					name="username"
					render={({ field }) => (
						<FormItem>
							<FormLabel>Username</FormLabel>
							<FormControl>
								<Input placeholder="Enter username" {...field} />
							</FormControl>
							<FormDescription>
								This is your public display name.
							</FormDescription>
							<FormMessage />
						</FormItem>
					)}
				/>
				<FormField
					control={form.control}
					name="email"
					render={({ field }) => (
						<FormItem>
							<FormLabel>Email</FormLabel>
							<FormControl>
								<Input type="email" placeholder="Enter email" {...field} />
							</FormControl>
							<FormDescription>
								We'll never share your email with anyone else.
							</FormDescription>
							<FormMessage />
						</FormItem>
					)}
				/>
				<FormField
					control={form.control}
					name="type"
					render={({ field }) => (
						<FormItem>
							<FormLabel>Account Type</FormLabel>
							<Select onValueChange={field.onChange} defaultValue={field.value}>
								<FormControl>
									<SelectTrigger>
										<SelectValue placeholder="Select account type" />
									</SelectTrigger>
								</FormControl>
								<SelectContent>
									<SelectItem value="personal">Personal</SelectItem>
									<SelectItem value="business">Business</SelectItem>
								</SelectContent>
							</Select>
							<FormDescription>
								Select the type of account you want to create.
							</FormDescription>
							<FormMessage />
						</FormItem>
					)}
				/>
				<FormField
					control={form.control}
					name="marketing"
					render={({ field }) => (
						<FormItem className="flex flex-row items-start space-x-3 space-y-0">
							<FormControl>
								<Checkbox
									checked={field.value}
									onCheckedChange={field.onChange}
								/>
							</FormControl>
							<div className="space-y-1 leading-none">
								<FormLabel>Marketing emails</FormLabel>
								<FormDescription>
									Receive emails about new products, features, and more.
								</FormDescription>
							</div>
						</FormItem>
					)}
				/>
				<Button type="submit">Submit</Button>
			</form>
		</Form>
	);
};

export const Default: Story = {
	args: {
		children: <SignUpForm />,
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Get form elements
		const usernameInput = canvas.getByLabelText("Username");
		const emailInput = canvas.getByLabelText("Email");
		const submitButton = canvas.getByRole("button", { name: "Submit" });

		// Try submitting empty form to trigger validation
		await userEvent.click(submitButton);

		// Enter invalid username (too short)
		await userEvent.type(usernameInput, "a");
		await userEvent.tab();

		// Wait for validation message
		const errorMessage = await canvas.findByText(
			"Username must be at least 2 characters.",
			{},
			{ timeout: 1000 },
		);
		await expect(errorMessage).toBeVisible();

		// Fix the username
		await userEvent.clear(usernameInput);
		await userEvent.type(usernameInput, "validuser");
		await userEvent.tab();

		// Verify error message is gone
		await expect(errorMessage).not.toBeVisible();
	},
};

const loginSchema = z.object({
	email: z.string().email(),
	password: z.string().min(8),
	remember: z.boolean().default(false),
});

const LoginForm = () => {
	const form = useForm<z.infer<typeof loginSchema>>({
		resolver: zodResolver(loginSchema),
		defaultValues: {
			email: "",
			remember: false,
		},
	});

	function onSubmit(values: z.infer<typeof loginSchema>) {}

	return (
		<Form {...form}>
			<form
				onSubmit={form.handleSubmit(onSubmit)}
				className="space-y-6 w-[400px]"
			>
				<FormField
					control={form.control}
					name="email"
					render={({ field }) => (
						<FormItem>
							<FormLabel>Email</FormLabel>
							<FormControl>
								<Input type="email" placeholder="Enter email" {...field} />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>
				<FormField
					control={form.control}
					name="password"
					render={({ field }) => (
						<FormItem>
							<FormLabel>Password</FormLabel>
							<FormControl>
								<Input
									type="password"
									placeholder="Enter password"
									{...field}
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>
				<FormField
					control={form.control}
					name="remember"
					render={({ field }) => (
						<FormItem className="flex flex-row items-start space-x-3 space-y-0">
							<FormControl>
								<Checkbox
									checked={field.value}
									onCheckedChange={field.onChange}
								/>
							</FormControl>
							<div className="space-y-1 leading-none">
								<FormLabel>Remember me</FormLabel>
							</div>
						</FormItem>
					)}
				/>
				<Button type="submit" className="w-full">
					Sign in
				</Button>
			</form>
		</Form>
	);
};

export const Login: Story = {
	args: {
		children: <LoginForm />,
	},
};
