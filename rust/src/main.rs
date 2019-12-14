
use std::fmt;
use std::vec::Vec;
use core::borrow::Borrow;
use std::io::Read;
use std::ptr::null;

#[derive(Clone)]
pub struct Node {
    order: usize,
    count: usize,
    values: Vec<i32>,
    children: Vec<Option<Box<Node>>>
}
type LinkNode = Option<Box<Node>>;

fn get_mutable(link: &mut LinkNode) -> Option<&mut Node> {
    let mut node : Option<&mut Node> = link.as_mut().map(|node| &mut **node);
    return node;
}

enum BtreeStatus {
    Overflow,
    Underflow,
    Ok
}

impl Node {
    fn new(order:usize) -> Self {
        Node {
            order: order,
            count: 0,
            values: vec![0; order+1],
            children: vec![None; order+2]
        }
    }
    fn insertInto(&mut self, position :usize, value:i32) {
        if self.count > self.order {
            panic!("error horror, self.count: {}, self.order: {}", self.count, self.order);
        }
        let mut index = self.count;
        while index > position {
            self.children[index+1] = self.children[index].take();
            self.values[index] = self.values[index-1];
            index -= 1;
        }
        self.values[index] = value;
        self.count += 1;
    }

    fn recursiveInsert(&mut self, value:i32) -> BtreeStatus {
        let mut index = 0;
        while index < self.count && self.values[index] < value {
            index += 1;
        }
        if let Some(child) = get_mutable(&mut self.children[index]) {
            match child.recursiveInsert(value) {
                BtreeStatus::Overflow => {
                    println!("split this node");
                    splitNode(self, index);
                }
                _ => {}
            }
        } else {
            self.insertInto(index, value);
        }
        if self.count == self.order {
            return BtreeStatus::Overflow;
        }
        return BtreeStatus::Ok;
    }
}
fn splitNode(parent: &mut Node, position:usize) {
    if parent.children[position].is_none() {
        panic!("child node at position {} is None", position);
    }
    let mut ptr = get_mutable(&mut parent.children[position]).unwrap();
    let mid = ptr.order / 2;
    let midValue = ptr.values[mid];
    let mut left =  Option::from(Box::new(Node::new(ptr.order)));
    let mut right=  Option::from(Box::new(Node::new(ptr.order)));
    let mut left_ptr = get_mutable(&mut left).unwrap();
    let mut right_ptr = get_mutable(&mut right).unwrap();

    let mut index = 0;
    while index < mid {
        if ptr.children[index].is_some() {
            left_ptr.children[left_ptr.count] = ptr.children[index].take();
        }
        left_ptr.insertInto(left_ptr.count, ptr.values[index]);
        index += 1;
    }
    if ptr.children[index].is_some() {
        left_ptr.children[left_ptr.count] = ptr.children[index].take();
    }
    index += 1;
    while index < ptr.count {
        if ptr.children[index].is_some() {
            right_ptr.children[right_ptr.count] = ptr.children[index].take();
        }
        right_ptr.insertInto(right_ptr.count, ptr.values[index]);
        index += 1;
    }
    if ptr.children[index].is_some() {
        right_ptr.children[right_ptr.count] = ptr.children[index].take();
    }

    println!("\t left node size: {}", left_ptr.count);
    println!("\t right node size: {}", right_ptr.count);


    parent.insertInto(position, midValue);
    println!("\t parent node size: {}", parent.count);

    parent.children[position] = left.take();
    parent.children[position+1] = right.take();

}

fn splitRoot(ptr: &mut Node) {

    let mid = ptr.count / 2;
    let midValue = ptr.values[mid];
    let mut left =  Option::from(Box::new(Node::new(ptr.order)));
    let mut right=  Option::from(Box::new(Node::new(ptr.order)));
    let mut left_ptr = get_mutable(&mut left).unwrap();
    let mut right_ptr = get_mutable(&mut right).unwrap();

    let mut index = 0;
    while index < mid {
        if ptr.children[index].is_some() {
            left_ptr.children[left_ptr.count] = ptr.children[index].take();
        }
        left_ptr.insertInto(left_ptr.count, ptr.values[index]);
        index += 1;
    }
    if ptr.children[index].is_some() {
        left_ptr.children[left_ptr.count] = ptr.children[index].take();
    }
    index += 1;
    while index < ptr.count {
        if ptr.children[index].is_some() {
            right_ptr.children[right_ptr.count] = ptr.children[index].take();
        }
        right_ptr.insertInto(right_ptr.count, ptr.values[index]);
        index += 1;
    }
    if ptr.children[index].is_some() {
        right_ptr.children[right_ptr.count] = ptr.children[index].take();
    }

    println!("\t rootSplit: left node size: {}", left_ptr.count);
    println!("\t rootSplit: right node size: {}", right_ptr.count);

    ptr.count = 1;
    println!("\t rootSplit: parent node size: {}", ptr.count);

    ptr.values[0] = midValue;
    ptr.children[0] = left.take();
    ptr.children[1] = right.take();


}

pub struct Btree{
    root: LinkNode,
    height: usize
}


impl Btree {
    pub fn new(order:usize) -> Self {
        Btree{
            root: Option::from(Box::new(Node::new(order))),
            height: 0
        }
    }
    pub fn insert(&mut self, value:i32) {
        if let Some(root ) = get_mutable(&mut self.root) {
            let status = root.recursiveInsert(value);
            match status {
                BtreeStatus::Overflow => {
                    println!("split this root");
                    splitRoot(root);
                }
                _ => {}
            }
            println!("root size: {}", root.count);
        }
    }
    pub fn print(&mut self) {
        recursivePrint(get_mutable(&mut self.root), 0)
    }
}

fn recursivePrint(ptr: Option<&mut Node>, level:usize) {
    if let Some(node) = ptr {
        let mut index = node.count;
        while index > 0 {
            recursivePrint(get_mutable(&mut node.children[index]), level + 1);
            for x in 0..level*5 {
                print!(".");
            }
            println!("{}, ", node.values[index - 1]);
            index -= 1;
        }
        recursivePrint(get_mutable(&mut node.children[index]), level + 1);
    }
}

fn main() {

    let mut btree : Btree = Btree::new(3);
    for n in 1..=100 {
        println!("inserting: {} \n============================================", n);

        btree.insert(n);
        btree.print();
        println!("============================================");
    }
    println!("Hello, world!");
}
