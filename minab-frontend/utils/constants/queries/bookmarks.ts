// We are calling Go API directly for bookmark/unbookmark instead of Hasura actions/mutations.
// These placeholders remain so imports don't break; the actual calls are done via $fetch in pages.
export const BookmarkMutation = gql`query DummyBookmarkPlaceholder { __typename }`;
export const UnBookmarkMutation = gql`query DummyUnBookmarkPlaceholder { __typename }`;

export const GetBookmarksQuery = gql`
    query GetBookmarks(
        $user_id: Int!,
        $skip: Int = 0,
        $take: Int = 6,
    ) {       
        bookmarks_aggregate {
            aggregate {
                count
            }
        }
        bookmarks(
            where: {
            user_id: {
                _eq: $user_id
            }
            },
            offset: $skip,
            limit: $take
        ) {
            event {
                id
                title
                category {
                    name
                }
                location {
                    city
                    venue
                }
                thumbnail
                start_date
                end_date
            }  
        }
    }
`;