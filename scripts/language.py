import json
import os

rootPath = '../app/i18n/'

def insert(keys, value, parent):
    clonedKeys = list(keys);
    if len(clonedKeys) == 1:
        parent[clonedKeys[0]] = value
    else:
        newKey = clonedKeys.pop(0)
        newParent = parent.get(newKey, {})
        parent[newKey] = newParent
        insert(clonedKeys, value, newParent)
    return

for fn in os.listdir(rootPath + 'raw'):
    if fn[0] == '.':
        continue

    print('Migration of file:', fn)

    with open(rootPath + 'raw/' + fn, 'r') as jsonFile:
        data = json.load(jsonFile)

    migration = {}
    for key in data.keys():
        print('  Migrating', key)
        keys = []
        tempKey = key
        index = tempKey.find('.') + 1
        while index > 0:
            index = tempKey.find('.') + 1
            if index != 0:
                item = tempKey[:index-1]
            else:
                item = tempKey
            tempKey = tempKey[index:]
            keys.append(item)
        insert(keys, data[key], migration)

    with open(rootPath + fn, 'w+') as outputFile:
        json.dump(migration, outputFile)

    print('Done')
