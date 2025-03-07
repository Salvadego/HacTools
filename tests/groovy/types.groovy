typeManager = de.hybris.platform.jalo.type.TypeManager.getInstance()

(typeManager.getType('Product').allSubTypes + typeManager.getType('Product'))
.stream()
.map { itemType ->
    typeManager.getType(itemType.code).allSuperTypes
    .stream()
    .map { superType -> superType.code }
    .collect()
    .reverse()
    .join('/') + '/' + itemType.code
}
.collect()
.join('\n')
