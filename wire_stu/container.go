package main

import (
	"fmt"
	"reflect"
)

const (
	InjectTag = "inject"
)

type Container struct {
	beanConstructorMap map[string]func() any
	beanMap            map[string]any
}

var container *Container

func GetContainer() *Container {
	if container == nil {
		container = &Container{
			beanConstructorMap: make(map[string]func() any),
			beanMap:            make(map[string]any),
		}
	}
	return container
}

func (c *Container) RegisterBeanFactory(beanName string, beanConstruct func() any) {
	if _, ok := c.beanConstructorMap[beanName]; ok {
		panic("beanName already registered")
	}

	c.beanConstructorMap[beanName] = beanConstruct
}

func (c *Container) GetBean(beanName string) any {
	if bean, ok := c.beanMap[beanName]; ok {
		return bean
	}

	if beanConstruct, ok := c.beanConstructorMap[beanName]; ok {
		fmt.Printf("beanConstruct: %+v\n", beanConstruct)
		bean := beanConstruct()
		rv := reflect.ValueOf(bean)
		fmt.Printf("rv: %#v\n", rv)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		rt := rv.Type()
		fmt.Printf("rt: %#v\n", rt)

		for i := 0; i < rv.NumField(); i++ {
			injectTag := rt.Field(i).Tag.Get(InjectTag)
			fmt.Printf("injectTag: %#v\n", injectTag)
			if injectTag == "" {
				continue
			}
			injectBean := c.GetBean(injectTag)

			c.beanMap[injectTag] = injectBean
			rv.Field(i).Set(reflect.ValueOf(injectBean))
		}

		c.beanMap[beanName] = bean
		fmt.Printf("bean: %+v\n", bean)
		return bean
	}

	panic("beanName not registered")
}
