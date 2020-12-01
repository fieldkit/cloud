<template>
    <section class="container">
        <header> Notes & Comments </header>
        <form @submit.prevent="postComment" class="new-comment">
            <UserPhoto :user="user"> </UserPhoto>
            <input type="text" :placeholder="getNewCommentPlaceholder()" v-model="body">
            <span v-if="body"> Post </span>
        </form>
        <div class="list" v-if="posts.length > 0">
            <div class="subheader">
                <span class="comments-counter"> {{posts.length}} comments </span>
            </div>
            <div class="item" v-for="post in posts" v-bind:key="post.id">
                <UserPhoto :user="post.author"> </UserPhoto>
                <div class="flex column">
                    <span class="author"> {{post.author.name}} </span>
                    <span class="body"> {{post.body}} </span>
                    <div class="actions">

                    </div>
                </div>
            </div>
        </div>
    </section>
</template>

<script>
import CommonComponents from "@/views/shared";

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
    padding: 0 20px;
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

    input {
        height: 45px;
        padding: 14px 612px 12px 13px;
        border-radius: 2px;
        border: solid 1px $color-border;
        outline: none;
        width: 100%;
        font-weight: 500;
    }

    span {
        @include position(absolute, 50% 20px null null);
        transform: translateY(-50%);
        font-weight: 900;
    }

    img {
        margin-top: 0;
    }
}

.author {
    font-size: 16px;
    font-weight: 500;
    margin-bottom: 5px;
}

.body {
    max-width: 550px;
}

.item {
    @include flex(center);

    > div {
        @include flex();

        &.column {
            flex-direction: column;
        }
    }
}

</style>
