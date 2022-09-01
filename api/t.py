from genericpath import isdir
from os import listdir
import re
import os.path

EXPR = r"func \(r \*Resolver\) .*\(ctx context.Context, ([.\[]*)\) \(.*\) {"

text = """
func (r *Resolver) CheckoutLineDelete(ctx context.Context, checkoutID *string, lineID *string, token *uuid.UUID) (*gqlmodel.CheckoutLineDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CheckoutLinesAdd(ctx context.Context, checkoutID *string, lines []*gqlmodel.CheckoutLineInput, token *uuid.UUID) (*gqlmodel.CheckoutLinesAdd, error) {
	panic(fmt.Errorf("not implemented"))
}
"""

foundArgs = re.findall(EXPR, text)
for m in foundArgs:
    print(m)


# res = re.sub(EXPR, "--------", text)

# print(res)

# excludes = []


# def main():
#     items = sorted(listdir("."))
#     for item in items:
#         if not item.endswith(".resolvers.go") or item in excludes:
#             continue

#         existingFileRead = open(item, 'r')
#         existingContent = existingFileRead.read()
#         existingFileRead.close()

#         foundArgs: list[str] = re.findall(EXPR, existingContent)
#         if len(foundArgs) > 0:
#             splits = foundArgs[0].split(",")
#             replaceText = """args struct {"""

#             for split in splits:
#                 replaceText += "\n  {}".format(split.strip())

#             replaceText += "\n}"

#             # writeFile = open(item, 'w')
#             finalText = re.sub(EXPR, replaceText, existingContent)

#             print(finalText)
#             # writeFile.close()


# if __name__ == "__main__":
#     main()
