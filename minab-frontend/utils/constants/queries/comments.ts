export const SendCommentMutation = gql`
    mutation SendComment(
        $name: String!,
        $email: String!,
        $message: String
    ) {
        insert_comments(objects: {
            name: $name,
            email: $email,
            message: $message
        }) {
            affected_rows
        }
    }
`;
export const TOGGLE_BOOKMARK_MUTATION = gql`
  mutation ToggleBookmark($event_id: Int!, $user_id: Int!) {
    insert_event_bookmarks_one(object: {event_id: $event_id, user_id: $user_id}) {
      event_id
    }
  }
`;

export const DELETE_BOOKMARK_MUTATION = gql`
  mutation DeleteBookmark($event_id: Int!, $user_id: Int!) {
    delete_event_bookmarks_by_pk(event_id: $event_id, user_id: $user_id) {
      event_id
    }
  }
`;