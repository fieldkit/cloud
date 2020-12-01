<template>
    <section class="container">
        <header> Notes & Comments </header>
        <form @submit.prevent="postComment" class="new-comment">
            <UserPhoto :user="user"> </UserPhoto>
            <input type="text" :placeholder="getNewCommentPlaceholder()" v-model="body">
            <span class="new-comment-submit" v-if="body"> Post </span>
        </form>
        <div class="list" v-if="posts.length > 0">
            <div class="subheader">
                <span class="comments-counter"> {{posts.length}} comments </span>
            </div>
            <div class="item" v-for="post in posts" v-bind:key="post.id">
                <UserPhoto :user="post.author"> </UserPhoto>
                <div class="flex column">
                    <span class="timestamp"> {{formatTimestamp(post.createdAt)}} </span>
                    <span class="author"> {{post.author.name}} </span>
                    <span class="body"> {{post.body}} </span>
                    <div class="actions">
                        <span>
                            <img src="@/assets/icon-reply.svg" alt="">
                            Reply
                        </span>
                    </div>
                </div>
            </div>
        </div>
    </section>
</template>

<script>
import CommonComponents from "@/views/shared";
import moment from "moment";

export default {
    name: "Comments",
    components: {
        ...CommonComponents,
    },
    props: {
        type: {
            required: true,
        },
        user: {
            required: true,
        },
        project: {
            required: true,
        },
    },
    methods: {
        getNewCommentPlaceholder() {
            console.log("project", this.project);

            if (this.type === "project") {
                return "Comment on Project";
            } else {
                return "Write a comment about this Data View";
            }
        },
        postComment() {
            const comment = {
                projectID: this.project.id,
                body: this.body,
            };
            return this.$services.api.postComment(comment).then(data => {
                console.log("insert response", data);
            });
        },
        formatTimestamp(timestamp) {
            return moment(timestamp).fromNow();
        },
    },
    mounted() {
        console.log("yo");
        return this.$services.api.getComments(this.project.id).then((data) => {
            this.posts = data.posts;
            console.log("this comments", this.posts);
        });
    },
    data() {
        return {
            body: null,
            posts: [],
        };
    },
};
</script>

<style lang="scss" scoped>
@import "../../scss/global";

* {
    box-sizing: border-box;
    font-size: 14px;
}
.hide {
    display: none;
}

.container {
    margin-top: 20px;
    padding: 0 20px 30px;
    background: #fff;
    border-radius: 1px;
    border: 1px solid $color-border;
    box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.05);
}

header {
    @include flex(center, space-between);
    height: 52px;
    border-bottom: 1px solid $color-border;
    font-size: 20px;
    font-weight: 500;
}

.subheader {
    @include flex(center, space-between);
    border-top: 1px solid $color-border;
    border-bottom: 1px solid $color-border;
    padding: 15px 0;
}

.new-comment {
    padding: 22px 0;
    position: relative;
    @include flex(center);

    img {
        margin-top: 0;
        width: 46px;
        height: 46px;
    }

    input {
        height: 45px;
        padding: 14px 72px 12px 13px;
        border-radius: 2px;
        border: solid 1px $color-border;
        outline: none;
        width: 100%;
        font-weight: 500;

        &::placeholder {
            color: #cccdcf;
        }
    }

    &-submit {
        @include position(absolute, 50% 10px null null);
        @include flex(center);
        height: 45px;
        padding: 0 10px;
        transform: translateY(-50%);
        font-weight: 900;
    }
}

.comments-counter {
    font-family: $font-family-light;
}

.author {
    font-size: 16px;
    font-weight: 500;
    margin-bottom: 5px;
}

.body {
    max-width: 550px;
    font-family: $font-family-light;
}

.item {
    @include flex(flex-start);
    padding: 15px 0 0;

    &::v-deep img {
        margin-top: 0;
        width: 30px;
        height: 30px;
    }

    > div {
        @include flex();
    }
}

.column {
    width: 100%;
    flex-direction: column;
    border-bottom: 1px solid $color-border;
    position: relative;
}

.actions {
    margin: 15px 0;

    span {
        font-weight: 500;
        margin-right: 20px;
    }

    img {
        width: 14px;
        height: 11px;
    }
}

.timestamp {
    @include position(absolute, 0 0 null null);
    font-family: $font-family-light;
}

</style>
