import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import { Card, CardContent } from "../../components/card.js";
import {
	Carousel,
	CarouselContent,
	CarouselItem,
	CarouselNext,
	CarouselPrevious,
} from "../../components/carousel.js";

type CarouselProps = React.ComponentProps<typeof Carousel>;

const meta: Meta<CarouselProps> = {
	title: "Design System/Carousel",
	component: Carousel,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Carousel>;

const slides = [
	{ id: "slide-1", content: "1" },
	{ id: "slide-2", content: "2" },
	{ id: "slide-3", content: "3" },
	{ id: "slide-4", content: "4" },
	{ id: "slide-5", content: "5" },
];

const images = [
	{
		id: "img-1",
		src: "https://images.unsplash.com/photo-1615247001958-f4bc92fa6a4a?w=300&dpr=2&q=80",
		alt: "Food dish with garnish",
	},
	{
		id: "img-2",
		src: "https://images.unsplash.com/photo-1513104890138-7c749659a591?w=300&dpr=2&q=80",
		alt: "Pizza with fresh ingredients",
	},
	{
		id: "img-3",
		src: "https://images.unsplash.com/photo-1618219740975-d40978bb7378?w=300&dpr=2&q=80",
		alt: "Pasta dish with sauce",
	},
];

export const Default: Story = {
	render: () => (
		<div className="w-full max-w-sm">
			<Carousel>
				<CarouselContent>
					{slides.map((slide) => (
						<CarouselItem key={slide.id}>
							<Card>
								<CardContent className="flex aspect-square items-center justify-center p-6">
									<span className="text-4xl font-semibold">
										{slide.content}
									</span>
								</CardContent>
							</Card>
						</CarouselItem>
					))}
				</CarouselContent>
				<CarouselPrevious />
				<CarouselNext />
			</Carousel>
		</div>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const nextButton = canvas.getByRole("button", { name: "Next slide" });
		const prevButton = canvas.getByRole("button", { name: "Previous slide" });

		await expect(nextButton).toBeInTheDocument();
		await expect(prevButton).toBeInTheDocument();
		await expect(prevButton).toBeDisabled();

		await new Promise((resolve) => setTimeout(resolve, 500));
		await userEvent.click(nextButton);
		await expect(prevButton).toBeEnabled();
	},
};

export const Vertical: Story = {
	render: () => (
		<div className="w-full max-w-sm">
			<Carousel orientation="vertical" className="h-[300px]">
				<CarouselContent>
					{slides.map((slide) => (
						<CarouselItem key={slide.id}>
							<Card>
								<CardContent className="flex aspect-square items-center justify-center p-6">
									<span className="text-4xl font-semibold">
										{slide.content}
									</span>
								</CardContent>
							</Card>
						</CarouselItem>
					))}
				</CarouselContent>
				<CarouselPrevious />
				<CarouselNext />
			</Carousel>
		</div>
	),
};

export const WithImages: Story = {
	render: () => (
		<div className="w-full max-w-sm">
			<Carousel>
				<CarouselContent>
					{images.map((image) => (
						<CarouselItem key={image.id}>
							<Card>
								<CardContent className="flex aspect-square items-center justify-center p-6">
									<img
										src={image.src}
										alt={image.alt}
										className="w-full h-full object-cover rounded-lg"
									/>
								</CardContent>
							</Card>
						</CarouselItem>
					))}
				</CarouselContent>
				<CarouselPrevious />
				<CarouselNext />
			</Carousel>
		</div>
	),
};

export const AutoPlay: Story = {
	render: () => (
		<div className="w-full max-w-sm">
			<Carousel
				opts={{
					align: "start",
					loop: true,
				}}
			>
				<CarouselContent>
					{slides.map((slide) => (
						<CarouselItem key={slide.id}>
							<Card>
								<CardContent className="flex aspect-square items-center justify-center p-6">
									<span className="text-4xl font-semibold">
										{slide.content}
									</span>
								</CardContent>
							</Card>
						</CarouselItem>
					))}
				</CarouselContent>
				<CarouselPrevious />
				<CarouselNext />
			</Carousel>
		</div>
	),
};
