import gql from "graphql-tag";

export const uploadFileMutation = gql`
  mutation UploadFile($name: String!, $base64: String!) {
    uploadFile(name: $name, base64: $base64) {
      url
    }
  }
`;