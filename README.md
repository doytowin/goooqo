[![Sonar Stats](https://sonarcloud.io/api/project_badges/measure?project=win.doyto.goooqo&metric=alert_status)](https://sonarcloud.io/dashboard?id=win.doyto.goooqo)
[![Code Lines](https://sonarcloud.io/api/project_badges/measure?project=win.doyto.goooqo&metric=ncloc)](https://sonarcloud.io/component_measures?id=win.doyto.goooqo&metric=ncloc)
[![Coverage Status](https://sonarcloud.io/api/project_badges/measure?project=win.doyto.goooqo&metric=coverage)](https://sonarcloud.io/component_measures?id=win.doyto.goooqo&metric=coverage)

GoooQo - An OQM Implementation That Can Automatically Build SQL Statements from Objects
---

## Introduction to OQM

The biggest difference between OQM technology and traditional ORM (object-relational mapping) technology is that OQM proposes to build CRUD statements directly through objects.

The core function of OQM is to build a query clause through a query object, which is also the origin of the Q in the name of OQM.

Another significant discovery in OQM technology is that the field names in query objects and the conditions in query clauses can be converted interchangeably.

In this way, we only need to create an entity object and a query object to build CRUD statements. 
The entity object is used to determine the table name and the column names, 
and the instance of the query object is used to control the construction of the query clause.

## Introduction to GoooQo

GoooQo is an OQM implementation that can automatically build SQL statements from objects.

The first three Os in the name `GoooQo` stands for the three major object concepts in the OQM technique:

- `Entity Object` is used to map the static part in the SQL statements, such as table name and column names;
- `Query Object` is used to map the dynamic part in the SQL statements, such as filter conditions, pagination, and sorting;
- `View Object` is used to map the static part in the complex query statements, such as table names, column names, nested views, and group-by columns.

Where `Qo` represents `Query Object`, which is the core concept in the OSM technique.

Check this [article](https://blog.doyto.win/post/introduction-to-goooqo-en/) for more details. 

Check this [demo](https://github.com/doytowin/goooqo-demo) to take a quick tour.

Check our [wiki](https://github.com/doytowin/goooqo/wiki) for the incoming documentations.

> This is currently an experimental project and is not suitable for production usage.
