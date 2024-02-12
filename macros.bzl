

def buildium_report(name, config, input, balances, rows):
    native.genrule(
        name = "report",
        srcs = [
            config, input
        ],
        outs = [
            balances, rows
        ],
        tools = [
            "@fintools_public//cmd/buildium-csv-read",
        ],
        cmd = """
          $(location @fintools_public//cmd/buildium-csv-read) \
            --cfg=$(location {config}) \
            --csv=$(location {input}) \
            --bal=$(location {balances}) \
            --rows=$(location {rows})
        """.format(config=config, input=input, balances=balances, rows=rows),
    )
