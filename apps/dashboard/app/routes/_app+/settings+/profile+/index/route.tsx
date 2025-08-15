import { api } from "@autopilot/api";
import { Button } from "@autopilot/ui/components/button";
import { Input } from "@autopilot/ui/components/input";
import { Label } from "@autopilot/ui/components/label";
import { zodResolver } from "@autopilot/ui/lib/hook-form-zod-resolver";
import { useForm } from "@autopilot/ui/lib/react-hook-form";
import { toast } from "@autopilot/ui/lib/sonner";
import { z } from "@autopilot/ui/lib/zod";
import { useTranslation } from "react-i18next";
import { useIdentity } from "@/components/identity-provider";
import type { RouteHandle } from "@/routes/_app+/_layout";

interface UpdateProfileData {
	name: string;
	image: FileList;
}

export const handle = {
	breadcrumb: "common:breadcrumbs.settings.profile.index",
} satisfies RouteHandle;

export default function Component() {
	const { t } = useTranslation(["common"]);
	const { user, refreshUser } = useIdentity();
	const { mutateAsync: editProfileMutate } = api.v1.useMutation(
		"put",
		"/v1/users/{id}",
		{
			onSuccess: refreshUser,
		},
	);
	const { mutateAsync: uploadImageMutate } = api.v1.useMutation(
		"post",
		"/v1/users/{id}/image",
		{
			onSuccess: refreshUser,
		},
	);

	const updateProfileSchema = z.object({
		name: z.string().min(2, t("common:settings.profile.errors.nameMinLength")),
		image: z.instanceof(FileList),
	});

	const {
		register,
		handleSubmit,
		formState: { errors },
	} = useForm<UpdateProfileData>({
		resolver: zodResolver(updateProfileSchema),
		defaultValues: {
			name: user?.name || "",
		},
	});

	const onSubmit = handleSubmit(async (data) => {
		try {
			const img = data.image.item(0);
			if (img) {
				const bytes = await img.bytes();
				await uploadImageMutate({
					params: {
						header: { "X-File-Name": img.name },
						path: { id: "@me" },
					},
					body: "",
					bodySerializer: () => {
						return bytes;
					},
				});
			}
			if (data.name !== user?.name) {
				await editProfileMutate({
					params: {
						path: { id: "@me" },
					},
					body: {
						name: data.name,
					},
				});
			}
			toast.success(t("common:settings.profile.success"));
		} catch (error) {
			toast.error(t("common:settings.profile.error"));
		}
	});

	return (
		<div className="space-y-8">
			<div className="pb-8">
				<div className="grid grid-cols-1 gap-x-8 gap-y-4 lg:grid-cols-2">
					<div>
						<h3 className="text-lg font-medium">
							{t("common:settings.profile.title")}
						</h3>

						<p className="text-sm text-muted-foreground">
							{t("common:settings.profile.description")}
						</p>
					</div>

					<form onSubmit={onSubmit} className="space-y-4">
						<div className="space-y-2">
							<Label htmlFor="email">
								{t("common:settings.profile.email")}
							</Label>

							<Input
								type="email"
								value={user?.email}
								disabled
								className="bg-muted"
							/>
						</div>

						<div className="space-y-2">
							<Label htmlFor="name">{t("common:settings.profile.name")}</Label>

							<Input
								{...register("name")}
								className={errors.name ? "ring-2 ring-destructive" : ""}
							/>
							{errors.name?.message && (
								<p className="text-sm text-destructive">
									{errors.name.message}
								</p>
							)}
						</div>

						<div className="w-64 space-y-2">
							<Label htmlFor="image">
								{t("common:settings.profile.picture")}
							</Label>

							<Input
								{...register("image")}
								className="pt-1.5"
								type="file"
								accept=".png,.jpg,.jpeg,.webp"
							/>
						</div>

						<div className="flex justify-end mt-8">
							<Button type="submit">
								{t("common:settings.profile.updateProfile")}
							</Button>
						</div>
					</form>
				</div>
			</div>
		</div>
	);
}
