import Config from "@/secrets";
import { DisplayProject } from "@/store";

export function twitterCardMeta(project: DisplayProject) {
    return [
        { name: "twitter:card", content: "summary_large_image" },
        { name: "twitter:site", content: "@FieldKitOrg" },
        { name: "twitter:title", content: project.name },
        { name: "twitter:description", content: project.project.description },
        { name: "twitter:image", content: Config.baseUrl + project.project.photo },
        { name: "twitter:image:alt", content: project.project.description }, // TODO It would be nice if we could allow users to set this. Maybe warn them description is used here?
    ];
}
