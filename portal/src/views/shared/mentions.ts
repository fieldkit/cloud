import DefaultProfilePhoto from "../../assets/profile-image-placeholder.svg";

export interface TributeOptions {
    values: (text, cb) => void;
    lookup: string;
    fillAttr: string;
    menuItemTemplate: unknown;
}

export function makeTributeOptions(config, services) {
    const baseUrl = config.baseUrl;
    return {
        values: (text, cb) => {
            if (text.length > 0) {
                services.api.mentionables(text).then((mentionables) => {
                    console.log("mentionables", mentionables);
                    cb(mentionables.users);
                });
            } else {
                cb([]);
            }
        },
        lookup: "name",
        fillAttr: "mention",
        menuItemTemplate: function(item) {
            if (item.original.photo) {
                return '<img class="tribute-user-icon" src="' + baseUrl + item.original.photo.url + '">' + item.string;
            } else {
                return '<img class="tribute-user-icon" src="' + DefaultProfilePhoto + '" />' + item.string;
            }
        },
        /*
        selectTemplate: function(item) {
            return (
                '<span contenteditable="false"><a href="http://zurb.com" target="_blank" title="' +
                item.original.email +
                '">' +
                item.original.value +
                "</a></span>"
            );
		},
		*/
    };
}
